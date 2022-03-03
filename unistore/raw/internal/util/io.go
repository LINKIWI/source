package util

import (
	"io"
)

// nopWriteCloser is an io.WriteCloser that wraps an io.Writer and provides a noop Close method.
type nopWriteCloser struct {
	io.Writer
}

// NopWriteCloser creates an io.WriteCloser from an io.Writer with a noop Close implementation.
// It is the io.Writer analog for io.NopCloser.
func NopWriteCloser(writer io.Writer) io.WriteCloser {
	return &nopWriteCloser{writer}
}

// Close is a noop.
func (nwc *nopWriteCloser) Close() error {
	return nil
}
