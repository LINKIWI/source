package mcrpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"lib.kevinlin.info/mcrpc/internal/server"
	"lib.kevinlin.info/mcrpc/protocol"
)

// Server is a memcache protocol server.
type Server struct {
	handler          Handler
	connReadTimeout  time.Duration
	connWriteTimeout time.Duration
	errorLog         *log.Logger
	warnLog          *log.Logger
	infoLog          *log.Logger
	debugLog         *log.Logger
}

// NewServer creates a new server backed by a Handler.
func NewServer(handler Handler, opts ...Option) (*Server, error) {
	s := &Server{handler: handler}

	for _, opt := range opts {
		s = opt(s)
	}

	return s, nil
}

// Serve starts the server on the specified listener.
func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.errorLog != nil {
				s.errorLog.Printf(
					"server: error accepting new connection on listener: err=%v",
					err,
				)
			}

			continue
		}

		if s.debugLog != nil {
			s.debugLog.Printf(
				"server: accepted client connection: laddr=%v raddr=%v",
				conn.LocalAddr(),
				conn.RemoteAddr(),
			)
		}

		conn = &server.TimeoutConn{
			ReadTimeout:  s.connReadTimeout,
			WriteTimeout: s.connWriteTimeout,
			Conn:         conn,
		}

		go s.handle(conn)
	}
}

// handle processes the client connection.
func (s *Server) handle(conn net.Conn) {
	reader := protocol.NewReader(conn)

	// Start a loop reading and processing requests from the client, in order to support
	// connection reuse across multiple requests.
	for {
		if s.debugLog != nil {
			s.debugLog.Printf(
				"server: read-evaluation loop init: conn=%v",
				conn.RemoteAddr(),
			)
		}

		terminate, err := s.dispatch(conn, reader)
		if err != nil && s.errorLog != nil {
			s.errorLog.Printf("server: error handling connection: err=%v", err)
		}

		if terminate {
			if s.debugLog != nil {
				s.debugLog.Printf(
					"server: terminating read-evaluation loop: conn=%v",
					conn.RemoteAddr(),
				)
			}

			conn.Close()
			return
		}
	}
}

// dispatch parses the request and invokes the underlying service handler.
func (s *Server) dispatch(conn net.Conn, reader *protocol.Reader) (bool, error) {
	var req protocol.Request
	var resp protocol.Response
	var hErr error

	buf, err := reader.ReadASCIICommand()
	if err == io.EOF {
		return true, nil
	} else if err != nil {
		return true, fmt.Errorf(
			"server: error buffering command for parse: buf=%#v err=%v",
			string(buf),
			err,
		)
	}

	req, err = protocol.NewASCIIParser().Parse(buf)
	if err == protocol.ErrInvalidParse {
		resp = &protocol.ClientError{Err: err}
		if _, cErr := conn.Write([]byte(resp.String())); cErr != nil {
			return true, cErr
		}
		return false, nil
	} else if err != nil {
		return true, fmt.Errorf(
			"server: non-recoverable error during command parsing: err=%v",
			err,
		)
	}

	if s.debugLog != nil {
		s.debugLog.Printf("server: parsed request: request=%#v", req)
	}

	ctx := context.WithValue(context.Background(), ClientConnContextKey, conn)
	ctx = context.WithValue(ctx, ClientRawRequestContextKey, buf)

	switch r := req.(type) {
	case *protocol.VersionRequest:
		resp, hErr = s.handler.Version(ctx, r)
	case *protocol.ShutdownRequest:
		resp, hErr = s.handler.Shutdown(ctx, r)
	case *protocol.FlushRequest:
		resp, hErr = s.handler.Flush(ctx, r)
	case *protocol.QuitRequest:
		resp, hErr = s.handler.Quit(ctx, r)
	case *protocol.StatsRequest:
		resp, hErr = s.handler.Stats(ctx, r)
	case *protocol.WatchRequest:
		resp, hErr = s.handler.Watch(ctx, r)
	case *protocol.TouchRequest:
		resp, hErr = s.handler.Touch(ctx, r)
	case *protocol.DeleteRequest:
		resp, hErr = s.handler.Delete(ctx, r)
	case *protocol.IncrRequest:
		resp, hErr = s.handler.Incr(ctx, r)
	case *protocol.DecrRequest:
		resp, hErr = s.handler.Decr(ctx, r)
	case *protocol.GetRequest:
		resp, hErr = s.handler.Get(ctx, r)
	case *protocol.GetsRequest:
		resp, hErr = s.handler.Gets(ctx, r)
	case *protocol.GatRequest:
		resp, hErr = s.handler.Gat(ctx, r)
	case *protocol.GatsRequest:
		resp, hErr = s.handler.Gats(ctx, r)
	case *protocol.SetRequest:
		resp, hErr = s.handler.Set(ctx, r)
	case *protocol.AddRequest:
		resp, hErr = s.handler.Add(ctx, r)
	case *protocol.ReplaceRequest:
		resp, hErr = s.handler.Replace(ctx, r)
	case *protocol.AppendRequest:
		resp, hErr = s.handler.Append(ctx, r)
	case *protocol.PrependRequest:
		resp, hErr = s.handler.Prepend(ctx, r)
	case *protocol.CasRequest:
		resp, hErr = s.handler.Cas(ctx, r)
	default:
		return false, fmt.Errorf("server: unsupported request type: request=%T", req)
	}

	if req.IsNoReply() {
		return false, nil
	}

	if hErr != nil {
		if s.debugLog != nil {
			s.debugLog.Printf("server: handler error response: err=%v", hErr)
		}

		if _, cErr := conn.Write([]byte(hErr.Error())); cErr != nil {
			return true, cErr
		}
	} else {
		if s.debugLog != nil {
			s.debugLog.Printf("server: handler response: response=%#v", resp)
		}

		if _, cErr := conn.Write([]byte(resp.String())); cErr != nil {
			return true, cErr
		}
	}

	return false, nil
}
