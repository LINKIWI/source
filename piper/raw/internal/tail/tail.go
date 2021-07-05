package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"piper/internal/config"
)

// State is an enum that describes the current state of the client within its finite state machine
// logic. Note that certain operations are only permitted when the client is in a particular state.
type State int

const (
	// Idle is a newly or partially instantiated Tail that does not yet have one or both of a
	// file descriptor on the log file open or an inotify watch registered on the log file.
	Idle State = iota
	// Consuming indicates that a file descriptor is opened and an inotify watch has been
	// successfully registered, and the client is either ready to send tail events or is
	// currently sending tail events to a consumer.
	Consuming
	// Errored indicates that the client has encountered an error. Consumption cannot proceed
	// in this state.
	Errored
)

var (
	// errNoData indicates that a read was attempted on a buffer with an incomplete log line,
	// and thus no new lines are available for consumption.
	errNoData = fmt.Errorf("tail: no line currently available due to incomplete data")
)

// Tail provides an abstraction for realtime consumption of a log file on disk.
//
// By default, on initialization, it opens a handle to a file and initializes a filesystem watch
// on the associated inode with inotify. Tail seeks to the end of the file and attempts to read a
// new line (defined as a sequence of bytes terminating in the newline character '\n') every time
// inotify reports a write event on the file.
//
// Internally, Tail is a state machine that transitions between Idle, Consuming, and Errored states.
// All errors encountered during consumption bring Tail into the Errored state, and the latest error
// is available via Error method. Most notably, io.EOF is raised when the consumption reads past the
// end of the file: this can happen if the log file was truncated (e.g. with a logrotate rule that
// uses the copytruncate directive). It is expected that the client handle such errors, which
// usually involves either re-creating the Tail instance or calling its Reset method.
type Tail struct {
	path    string
	file    *os.File
	reader  *bufio.Reader
	watcher *Watcher
	state   State
	opts    Options

	pending string
	err     error
}

// Options provides parameters for tuning the runtime behavior of a Tail.
type Options struct {
	// Delimiter is the character constant separating new lines in the log file.
	Delimiter byte
	// BufferLength is the maximum capacity of the message channel into which messages can be
	// inserted without consumption. This allows for temporal decoupling of file I/O and
	// downstream consumption; i.e., a new tailed line can be read from the file immediately
	// without waiting for the consumer to finish processing the previous line. However, no
	// matter the buffer length, messages will still be published for every append change
	// notification.
	BufferLength int
	// SeekPosition describes the file seek behavior to employ before consuming tail events.
	SeekPosition config.SeekPosition
	// Tags is an arbitrary key-value map of metadata to attach to every published message.
	Tags map[string]string `json:"tags"`
}

// Message is a structured event representing a new tailed line in the log file.
type Message struct {
	// Timestamp is the timestamp at which the line was written.
	Timestamp time.Time `json:"timestamp"`
	// SequenceID is an incrementing counter that indicates the number of lines read from the
	// file since the line associated with this message.
	SequenceID int `json:"sequence_id"`
	// Path is the path to the file that raised the associated inotify event.
	Path string `json:"path"`
	// Line is the string contents of the newly written line.
	Line string `json:"line"`
	// Tags is an arbitrary key-value map of metadata supplied by the client.
	Tags map[string]string `json:"tags"`
}

// New creates and attempts to initialize a Tail.
// A successful initialization will transition the Tail into a Consuming state; otherwise, an error
// is returned.
func New(path string, opts Options) (*Tail, error) {
	tail := &Tail{
		path:  path,
		state: Idle,
		opts:  opts,
	}

	if tail.opts.Delimiter == 0 {
		tail.opts.Delimiter = '\n'
	}

	if err := tail.Reset(); err != nil {
		return nil, err
	}

	return tail, nil
}

// Reset attempts to re-initialize a Tail by closing and reopening all necessary file handles and
// filesystem watches.
func (t *Tail) Reset() error {
	var err error

	t.file = nil
	t.watcher = nil
	t.reader = nil
	t.err = nil
	t.pending = ""
	t.state = Idle

	if err := t.Close(); err != nil {
		return err
	}

	t.file, err = os.OpenFile(t.path, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return err
	}

	stat, err := t.file.Stat()
	if err != nil {
		return err
	}

	if stat.Mode().IsRegular() {
		if t.opts.SeekPosition == config.SeekPositionEnd {
			_, err = t.file.Seek(0, io.SeekEnd)
			if err != nil {
				return err
			}
		}
	} else if stat.Mode()&os.ModeNamedPipe != os.ModeNamedPipe {
		return fmt.Errorf("tail: file is not a regular file nor a named pipe")
	}

	t.watcher, err = NewWatcher(t.path)
	if err != nil {
		return err
	}

	t.reader = bufio.NewReader(t.file)
	t.state = Consuming

	return nil
}

// Consume allows realtime consumption of tailed lines from a log file over an asynchronous channel.
// Starting consumption reads inotify events from the file watcher and extracts lines that end with
// the newline delimiter for transmission over the channel. The channel is closed when an error is
// encountered or when the inotify watch is removed.
func (t *Tail) Consume() (<-chan Message, error) {
	messages := make(chan Message, t.opts.BufferLength)
	sequenceID := 0

	if t.state != Consuming {
		return nil, fmt.Errorf(
			"tail: consumption is not permitted in the current state: state=%v",
			t.state,
		)
	}

	go func() {
		for {
			select {
			case event, ok := <-t.watcher.Events():
				if !ok {
					close(messages)
					return
				}

				switch event {
				case FsRemove, FsRename:
					t.err = fmt.Errorf("tail: file relocated")
					t.state = Errored
					close(messages)
					return

				case FsTruncate:
					t.err = fmt.Errorf("tail: file truncated")
					t.state = Errored
					close(messages)
					return

				case FsOpen:
					// When using config.SeekPositionStart, existing lines in
					// the file should be consumed. Thus, a file open event
					// should be treated identically to an append event.
					// When config.SeekPositionEnd, file open events should be
					// ignored since only newly appended lines since the end of
					// the file should be consumed.
					if t.opts.SeekPosition == config.SeekPositionEnd {
						break
					}
					fallthrough

				case FsAppend:
					lines, err := t.readLines()
					if err == errNoData || err == io.EOF {
						// No new log lines are currently available; wait
						// until the next write and try again.
						break
					} else if err != nil {
						t.err = err
						t.state = Errored
						close(messages)
						return
					}

					for _, line := range lines {
						messages <- Message{
							Timestamp:  time.Now(),
							SequenceID: sequenceID,
							Path:       t.path,
							Line:       line,
							Tags:       t.opts.Tags,
						}

						sequenceID++
					}
				}

			case err := <-t.watcher.Errors():
				if err != nil {
					t.err = err
					t.state = Errored
				}

				close(messages)
				return
			}
		}
	}()

	return messages, nil
}

// readLines attempts to read as many complete log lines as possible from the current seek position.
// errNoData is returned when zero log lines were read from the buffer while the seek position is
// not yet at the end of the file; this indicates an incomplete log line was written that may be
// processed on the next read attempt. All other errors, including io.EOF, should be considered
// unexpected and thus handled accordingly.
func (t *Tail) readLines() ([]string, error) {
	var lines []string

	line, err := t.reader.ReadString(t.opts.Delimiter)

	if err == io.EOF && len(line) == 0 {
		// Immediately encountering io.EOF without reading any bytes from the stream is a
		// fairly reliable indicator that the current seek position is invalid. This could,
		// for example, be due to truncation caused by log rotation, or a duplicated
		// filesystem notification event fired on file whose seek position is already at the
		// end of the file. In any case, the Tail client has no ability to determine the
		// reason for the error (and thus has no heuristics available to recover the
		// consumption) so the error is simply propagated downstream as-is.
		return lines, io.EOF
	} else if err == io.EOF && len(line) > 0 {
		// Encountering io.EOF while successfully reading back a nonzero number of bytes
		// from the stream indicates absence of the specified delimiter before reaching the
		// end of the file. This could, for example, arise from an incomplete or partially
		// written log line. It's likely that a future write will populate the delimiter, so
		// cache the data into a temporary buffer for consumption on subsequent reads that
		// do encounter a valid delimiter.
		t.pending += line
		return lines, errNoData
	} else if err != nil {
		// Unexpected error condition; propagate to downstream.
		return lines, err
	}

	// A line was read successfully. Any pending data from a previous read attempt should also
	// be finalized at this stage.
	lines = append(lines, t.pending+strings.TrimRight(line, string(t.opts.Delimiter)))
	t.pending = ""

	// However, there may be additional lines left in the buffer, accessible by recursion.
	// At this point, both an EOF and errNoData are graceful failure modes that indicate there
	// does not exist any more lines to read from the buffer.
	additionalLines, err := t.readLines()
	if err != nil && err != io.EOF && err != errNoData {
		return lines, err
	}

	return append(lines, additionalLines...), nil
}

// State provides the current state of the client.
func (t *Tail) State() State {
	return t.state
}

// Error reports the most recent error encountered by the client during consumption. if available.
// The error is guaranteed to be non-nil when the client is in the Errored state.
func (t *Tail) Error() error {
	return t.err
}

// Close attempts to gracefully release file descriptors currently held by the client.
func (t *Tail) Close() error {
	if t.watcher != nil {
		if err := t.watcher.Close(); err != nil {
			return err
		}
	}

	if t.file != nil {
		return t.file.Close()
	}

	return nil
}
