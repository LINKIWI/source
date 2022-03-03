package util

import (
	"io"
)

// InterceptedIOStream is an io.ReadWriteCloser that internally tracks the number of bytes read and
// written, and passes this value to an optional callback when the reader is closed.
type InterceptedIOStream struct {
	nread     int
	nwritten  int
	callbacks []func(nread int, nwritten int)
	io.Reader
	io.Writer
	io.Closer
}

// NewInterceptedIOStream wraps an existing I/O stream with size tracking across Read and Write
// invocations..
func NewInterceptedIOStream(stream interface{}) *InterceptedIOStream {
	is := &InterceptedIOStream{}

	if reader, ok := stream.(io.Reader); ok {
		is.Reader = reader
	}

	if writer, ok := stream.(io.Writer); ok {
		is.Writer = writer
	}

	if closer, ok := stream.(io.Closer); ok {
		is.Closer = closer
	}

	return is
}

// Read defers to the underlying io.Reader and increments an internal counter for the number of
// bytes read.
func (i *InterceptedIOStream) Read(p []byte) (n int, err error) {
	n, err = i.Reader.Read(p)
	i.nread += n
	return
}

// Write defers to the underlying io.Writer and increments an internal counter for the number of
// bytes written.
func (i *InterceptedIOStream) Write(p []byte) (n int, err error) {
	n, err = i.Writer.Write(p)
	i.nwritten += n
	return
}

// Close closes the underlying io.Closer and invokes the callback with the total number of
// bytes read from the reader.
func (i *InterceptedIOStream) Close() error {
	defer func() {
		for _, callback := range i.callbacks {
			callback(i.nread, i.nwritten)
		}
	}()

	return i.Closer.Close()
}

// RegisterCallback registers a callback to execute when the underlying stream is closed. The
// callback is invoked with the total number of bytes read and written to the io.Reader and
// io.Writer, respectively.
func (i *InterceptedIOStream) RegisterCallback(callback func(nread int, nwritten int)) {
	i.callbacks = append(i.callbacks, callback)
}

// BytesRead returns the current number of bytes read from the io.Reader.
func (i *InterceptedIOStream) BytesRead() int {
	return i.nread
}

// BytesWritten returns the current number of bytes written to the io.Writer.
func (i *InterceptedIOStream) BytesWritten() int {
	return i.nwritten
}
