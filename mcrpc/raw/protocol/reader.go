package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	// ErrMalformedRequest indicates that no protocol-compliant message could be extracted from
	// the reader stream.
	ErrMalformedRequest = errors.New("protocol: request message is malformed")
)

var (
	// crlf is a byte sequence conventionally marking the end of a protocol message.
	crlf = []byte("\r\n")

	// scanChunkSize is the size of each successive scanned chunk from the stream.
	scanChunkSize = 1

	// storageCommands is a slice of protocol command names that are responsible for storing
	// data. These commands prefix the only messages in the protocol that terminate on the
	// second CRLF sequence, rather than the first.
	storageCommands = [][]byte{
		[]byte("set"),
		[]byte("add"),
		[]byte("replace"),
		[]byte("append"),
		[]byte("prepend"),
		[]byte("cas"),
	}
)

// Reader abstracts over an io.Reader that acts as an input stream for memcache protocol commands.
type Reader struct {
	b *bufio.Reader
}

// NewReader creates a new protocol-aware reader from an io.Reader.
func NewReader(stream io.Reader) *Reader {
	return &Reader{bufio.NewReader(stream)}
}

// Read reads from the buffered io.Reader. Note that the buffering may cause the reader's current
// logical cursor position to be beyond its physical counterpart in the underlying io.Reader. In
// order to maintain a consistent view of the reader's true contents, the Read method transparently
// uses the buffered reader.
func (r *Reader) Read(p []byte) (int, error) {
	return r.b.Read(p)
}

// ReadASCIICommand reads a full ASCII command from the underlying io.Reader, exercising heuristics
// to determine how far into the stream data should be read.
func (r *Reader) ReadASCIICommand() ([]byte, error) {
	buf, err := r.readUntil(crlf)
	if err != nil {
		return buf, err
	}

	// Storage commands are the only protocol request messages that do conclude on the first
	// CRLF sequence. The read length following the first CRLF is context-dependent on the
	// claimed data size in the storage header. This routine performs a low-fidelity parse of
	// the storage header and, assuming it is valid, attempts to read ahead exactly the number
	// of bytes stated in the header for the data size.
	for _, cmd := range storageCommands {
		if bytes.HasPrefix(buf, cmd) {
			components := bytes.Fields(buf)
			if len(components) < 5 {
				return buf, fmt.Errorf(
					"protocol: bad storage header format: components=%d",
					len(components),
				)
			}

			size, err := strconv.Atoi(string(components[4]))
			if err != nil {
				return buf, fmt.Errorf(
					"protocol: error parsing size in storage header: err=%v",
					err,
				)
			}

			data := make([]byte, size)
			n, err := r.Read(data)
			buf = append(buf, data...)

			if n != size {
				return buf, fmt.Errorf(
					"protocol: cannot accommodate purported storage data "+
						"size from buffer: expect=%d actual=%d",
					size,
					n,
				)
			}

			tailer, err := r.readUntil(crlf)
			buf = append(buf, tailer...)

			return buf, err
		}
	}

	return buf, nil
}

// readUntil incrementally reads from the buffered stream in chunks until the delimiter sequence is
// located. Only those bytes that represent the fully-formed message are consumed from the buffered
// reader. If the sequence is located before EOF, a byte slice with the full message contents is
// returned; otherwise, an error indicates the command is malformed.
func (r *Reader) readUntil(sequence []byte) ([]byte, error) {
	var terminate bool

	buf := make([]byte, 0, 512)

	for !terminate {
		chunk, err := r.b.Peek(scanChunkSize)
		if err != nil {
			if err == io.EOF && len(buf) > 0 {
				// Terminate the chunk-read loop on the next iteration, relevant
				// only when a nonzero buffer has already been accumulated.
				terminate = true
			} else {
				return buf, err
			}
		}

		buf = append(buf, chunk...)

		idx := bytes.Index(buf, sequence)
		if idx != -1 {
			offset := idx + len(sequence)
			r.b.Discard(offset - len(buf) + len(chunk))
			return buf[:offset], nil
		}

		r.b.Discard(len(chunk))
	}

	return buf, ErrMalformedRequest
}
