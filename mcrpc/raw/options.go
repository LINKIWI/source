package mcrpc

import (
	"log"
	"time"
)

// Option is a type alias for a server option.
type Option func(s *Server) *Server

// WithConnReadTimeout sets a timeout for all reads on client connections.
func WithConnReadTimeout(timeout time.Duration) Option {
	return func(s *Server) *Server {
		s.connReadTimeout = timeout
		return s
	}
}

// WithConnWriteTimeout sets a timeout for all writes on client connections.
func WithConnWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) *Server {
		s.connWriteTimeout = timeout
		return s
	}
}

// WithErrorLog enables a custom error logger on the server.
func WithErrorLog(logger *log.Logger) Option {
	return func(s *Server) *Server {
		s.errorLog = logger
		return s
	}
}

// WithWarnLog enables a custom warn logger on the server.
func WithWarnLog(logger *log.Logger) Option {
	return func(s *Server) *Server {
		s.warnLog = logger
		return s
	}
}

// WithInfoLog enables a custom error logger on the server.
func WithInfoLog(logger *log.Logger) Option {
	return func(s *Server) *Server {
		s.infoLog = logger
		return s
	}
}

// WithDebugLog enables a custom debug logger on the server.
func WithDebugLog(logger *log.Logger) Option {
	return func(s *Server) *Server {
		s.debugLog = logger
		return s
	}
}
