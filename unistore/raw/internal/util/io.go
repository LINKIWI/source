package util

import (
	"context"
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
func (mwc *multiWriteCloser) Close() error {
	var errs []error

	for _, writer := range mwc.wc {
		if err := writer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("util: error closing multiple write closers: errs=%v", errs)
}

// contextReader is a context-aware io.Reader simulating cancelable, asynchronous I/O over an
// existing io.Reader.
type contextReader struct {
	ctx context.Context
	io.Reader
}

// NewContextReader creates a new io.Reader wrapper with the provided context and reader.
func NewContextReader(ctx context.Context, reader io.Reader) io.Reader {
	return &contextReader{
		ctx:    ctx,
		Reader: reader,
	}
}

// Read passes the read to the child io.Reader while respecting context cancellations.
func (cr *contextReader) Read(p []byte) (int, error) {
	type ioResult struct {
		n   int
		err error
	}

	if cr.ctx.Err() != nil {
		return 0, cr.ctx.Err()
	}

	result := make(chan ioResult)

	go func() {
		r := ioResult{}
		r.n, r.err = cr.Reader.Read(p)
		result <- r
	}()

	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	case r := <-result:
		return r.n, r.err
	}
}
