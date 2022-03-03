package log

import (
	"net/http"

	"go.uber.org/zap"
)

// RequestLogHandler wraps an http.Handler with transparent Zap request logging at debug level.
type RequestLogHandler struct {
	name string
	http.Handler
}

// NewRequestLogHandler creates a new stateless logging handler on top of an existing http.Handler.
func NewRequestLogHandler(name string, handler http.Handler) http.Handler {
	return &RequestLogHandler{
		name:    name,
		Handler: handler,
	}
}

// ServeHTTP creates a debug log trace about the incoming request and defers to the underlying
// handler's ServeHTTP.
func (l *RequestLogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	zap.L().Debug(
		"serving HTTP request",
		zap.String("listener", l.name),
		zap.String("host", req.Host),
		zap.String("url", req.RequestURI),
		zap.String("method", req.Method),
		zap.String("protocol", req.Proto),
		zap.String("user_agent", req.UserAgent()),
		zap.String("referer", req.Referer()),
		zap.String("remote_addr", req.RemoteAddr),
		zap.Bool("tls", req.TLS != nil),
	)

	l.Handler.ServeHTTP(rw, req)
}
