package tail

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

// FsEvent describes a filesystem event on a watched file.
type FsEvent int

const (
	// FsUnknown is a fallthrough, unknown event type.
	FsUnknown FsEvent = iota
	// FsOpen indicates that the file was opened for reading.
	FsOpen
	// FsAppend indicates that data was (likely) appended to the file.
	FsAppend
	// FsRemove indicates that the file was removed.
	FsRemove
	// FsRename indicates that the file was renamed.
	FsRename
	// FsTruncate indicates that the file was truncated.
	FsTruncate
)

// Watcher is an abstraction over an fsnotify.Watcher that wraps event processing with additional
// heuristics to expose more filesystem event types.
type Watcher struct {
	fsw    *fsnotify.Watcher
	events chan FsEvent
	errs   chan error
	path   string
	size   int64
}

// NewWatcher creates a new watcher on a path.
func NewWatcher(path string) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := fsw.Add(path); err != nil {
		return nil, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		fsw:    fsw,
		events: make(chan FsEvent),
		errs:   make(chan error),
		path:   path,
		size:   stat.Size(),
	}

	go watcher.watch()

	return watcher, nil
}

// Events is an accessor for the filesystem events channel.
func (w *Watcher) Events() <-chan FsEvent {
	return w.events
}

// Errors is an accessor for the errors channel.
func (w *Watcher) Errors() <-chan error {
	return w.errs
}

// watch starts an indefinite loop that consumes events from fsnotify.Watcher and re-emits them
// with a best-effort more granular classification.
func (w *Watcher) watch() {
	defer close(w.events)
	defer close(w.errs)

	w.events <- FsOpen

	for {
		select {
		case evt, ok := <-w.fsw.Events:
			if !ok {
				return
			}

			if evt.Op&fsnotify.Write > 0 {
				// inotify only knows that the file was written to; there's no
				// information whether the file was appended to, modified in place,
				// or truncated. The following logic attempts to make a guess using
				// some best-effort heuristics around changes in the file size.
				stat, err := os.Stat(w.path)
				if err != nil {
					w.errs <- err
					break
				}

				if stat.Size() < w.size {
					w.events <- FsTruncate
				} else {
					// Strictly speaking, insertion of bytes at any position in
					// the file other than the end would not be considered an
					// append operation even though it increases the file size.
					// However, this is relatively uncommon write pattern for a
					// log file, so just assume it was an append.
					w.events <- FsAppend
				}

				w.size = stat.Size()
			} else if evt.Op&fsnotify.Rename > 0 {
				w.events <- FsRename
			} else if evt.Op&fsnotify.Remove > 0 {
				w.events <- FsRemove
			}

		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}

			w.errs <- err
		}
	}
}

// Close closes the underlying filesystem watch.
func (w *Watcher) Close() error {
	return w.fsw.Close()
}
