package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	// DefaultDirectPipe is a DirectPipe with default parameters.
	DefaultDirectPipe = &DirectPipe{
		Name:     "default",
		DebugLog: log.New(io.Discard, "", 0),
		ErrorLog: log.New(io.Discard, "", 0),
	}
)

// Pipe describes an implementation that pipes data between source and destination connections.
type Pipe interface {
	// Do performs a bidirectional data copy between source and destination net.Conn instances
	// and returns the total number of bytes transferred.
	Do(dst net.Conn, src net.Conn) (srcSent int64, dstRecv int64, err error)
}

// DirectPipe is an implementation of a bidirectional Pipe between two net.Conn connections.
// In scenarios where it is supported, the pipe gracefully closes the read end of a connection to
// avoid attempting to pipe bytes when one side of the connection has closed.
type DirectPipe struct {
	// Name is a human-readable identifier for this pipe instance.
	Name string
	// ConnectionLifetime is the timeout value for the maximum total amount of time the pipe is
	// allowed to take. Connections whose total duration exceeds this value will be considered
	// errored. A zero-valued lifetime will allow connections to take unlimited time.
	ConnectionLifetime time.Duration
	// SourceReadTimeout is the socket read timeout for the source connection.
	SourceReadTimeout time.Duration
	// SourceWriteTimeout is the socket write timeout for the source connection.
	SourceWriteTimeout time.Duration
	// DestinationReadTimeout is the socket read timeout for the destination connection.
	DestinationReadTimeout time.Duration
	// DestinationWriteTimeout is the socket write timeout for the destination connection.
	DestinationWriteTimeout time.Duration
	// DebugLog is a standard library logger for debug messages.
	DebugLog *log.Logger
	// ErrorLog is a standard library logger for error messages
	ErrorLog *log.Logger
}

// Do performs a pipe between source and destination connections.
func (p *DirectPipe) Do(dst net.Conn, src net.Conn) (int64, int64, error) {
	var srcSent int64
	var dstRecv int64
	var err error

	errs := make(chan error, 2)
	srcClosed := make(chan bool)
	dstClosed := make(chan bool)

	srcIO := &timeoutConn{
		readTimeout:  p.SourceReadTimeout,
		writeTimeout: p.SourceWriteTimeout,
		Conn:         src,
	}
	dstIO := &timeoutConn{
		readTimeout:  p.DestinationReadTimeout,
		writeTimeout: p.DestinationWriteTimeout,
		Conn:         dst,
	}

	p.DebugLog.Printf(
		"proxy: pipe operation started: name=%s src=%v dst=%v",
		p.Name,
		src.RemoteAddr(),
		dst.RemoteAddr(),
	)

	go func() {
		defer func() {
			src.Close()
			srcClosed <- true
		}()

		srcSent, err = io.Copy(dstIO, srcIO)

		if err != nil {
			errs <- err
			p.ErrorLog.Printf(
				"proxy: error piping src -> dst: name=%s src=%v dst=%v err=%v",
				p.Name,
				src.RemoteAddr(),
				dst.RemoteAddr(),
				err,
			)
		} else {
			p.DebugLog.Printf(
				"proxy: copied bytes src -> dst: name=%s src=%v dst=%v bytes=%d",
				p.Name,
				src.RemoteAddr(),
				dst.RemoteAddr(),
				srcSent,
			)
		}
	}()

	go func() {
		defer func() {
			dst.Close()
			dstClosed <- true
		}()

		dstRecv, err = io.Copy(srcIO, dstIO)

		if err != nil {
			errs <- err
			p.ErrorLog.Printf(
				"proxy: error piping dst -> src: name=%s src=%v dst=%v err=%v",
				p.Name,
				src.RemoteAddr(),
				dst.RemoteAddr(),
				err,
			)
		} else {
			p.DebugLog.Printf(
				"proxy: copied bytes dst -> src: name=%s src=%v dst=%v bytes=%d",
				p.Name,
				src.RemoteAddr(),
				dst.RemoteAddr(),
				dstRecv,
			)
		}
	}()

	timeout := make(<-chan time.Time)
	if p.ConnectionLifetime > 0 {
		timeout = time.After(p.ConnectionLifetime)
	}

	select {
	case err := <-errs:
		return 0, 0, err
	case <-timeout:
		return 0, 0, fmt.Errorf(
			"proxy: total connection duration exceeded allowed lifetime: lifetime=%v",
			p.ConnectionLifetime,
		)
	case <-dstClosed:
		p.DebugLog.Printf(
			"proxy: dst connection closed: name=%s dst=%v",
			p.Name,
			dst.RemoteAddr(),
		)
		switch conn := src.(type) {
		case *net.TCPConn:
			conn.CloseRead()
		case *net.UnixConn:
			conn.CloseRead()
		case *tls.Conn:
			conn.CloseWrite()
		}
		<-srcClosed
	case <-srcClosed:
		p.DebugLog.Printf(
			"proxy: src connection closed: name=%s src=%v",
			p.Name,
			src.RemoteAddr(),
		)
		switch conn := dst.(type) {
		case *net.TCPConn:
			conn.CloseRead()
		case *net.UnixConn:
			conn.CloseRead()
		case *tls.Conn:
			conn.CloseWrite()
		}
		<-dstClosed
	}

	p.DebugLog.Printf(
		"proxy: pipe operation complete: name=%s src=%v dst=%v",
		p.Name,
		src.RemoteAddr(),
		dst.RemoteAddr(),
	)

	return srcSent, dstRecv, nil
}
