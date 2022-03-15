package util

import (
	"fmt"
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

// multiWriteCloser is an io.WriteCloser that abstracts over multiple io.WriteCloser instances. It
// is similar in spirit to the io.Writer returned by io.MultiWriter but also implements io.Closer.
type multiWriteCloser struct {
	wc []io.WriteCloser
	io.Writer
}

// NewMultiWriteCloser creates a new multiWriteCloser from zero or more io.WriteClosers.
func NewMultiWriteCloser(wc ...io.WriteCloser) io.WriteCloser {
	writers := make([]io.Writer, len(wc))

	for i, writer := range wc {
		writers[i] = writer.(io.Writer)
	}

	return &multiWriteCloser{
		wc:     wc,
		Writer: io.MultiWriter(writers...),
	}
}

// Close attempts to close all underlying write closers, collecting all errors as applicable.
func (m *multiWriteCloser) Close() error {
	var errs []error

	for _, writer := range m.wc {
		if err := writer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("util: error closing multiple write closers: errs=%v", errs)
}
