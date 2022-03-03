package server

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go.uber.org/zap"
	"lib.kevinlin.info/aperture/lib"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/util"
)

const (
	metricRestrictedListenerAccept          = "server.listener.restricted.accept"
	metricRestrictedListenerAuthorize       = "server.listener.restricted.authorize"
	metricEagerTLSListenerAccept            = "server.listener.eager_tls.accept"
	metricEagerTLSListenerHandshakeError    = "server.listener.eager_tls.handshake.error"
	metricEagerTLSListenerHandshakeDuration = "server.listener.eager_tls.handshake.duration"
)

// restrictedListener is a net.Listener that only accepts TCP connections from allowed remote IP
// addresses. Connections from disallowed source IPs are closed immediately.
type restrictedListener struct {
	name    string
	sources []config.CIDR
	net.Listener
}

// newRestrictedListener creates a new restricted listener.
func newRestrictedListener(name string, sources []config.CIDR, ln net.Listener) net.Listener {
	return &restrictedListener{
		name:     name,
		sources:  sources,
		Listener: ln,
	}
}

// Accept defers to the underlying net.Listener and additionally verifies that the accepted
// connection originates from a permitted IP block. Disallowed connections are closed immediately.
func (l *restrictedListener) Accept() (net.Conn, error) {
	tags := map[string]interface{}{"name": l.name}

	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, &util.Error{
			Namespace:    "server",
			Message:      "failed to accept on listener",
			StackedError: err,
		}
	}

	metrics.Client.Incr(metricRestrictedListenerAccept, tags)

	addr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		zap.L().Warn(
			"restricted listener rejecting non-TCP accepted connection",
			zap.Stringer("addr", conn.RemoteAddr()),
		)
		metrics.Client.Incr(
			metricRestrictedListenerAuthorize,
			util.MergeMaps(tags, map[string]interface{}{"action": "reject"}),
		)

		conn.Close()
		return conn, nil
	}

	for _, source := range l.sources {
		if source.Contains(addr.IP) {
			metrics.Client.Incr(
				metricRestrictedListenerAuthorize,
				util.MergeMaps(tags, map[string]interface{}{"action": "accept"}),
			)
			return conn, nil
		}
	}

	zap.L().Warn(
		"restricted listener rejecting unauthorized accepted connection",
		zap.Stringer("addr", conn.RemoteAddr()),
	)
	metrics.Client.Incr(
		metricRestrictedListenerAuthorize,
		util.MergeMaps(tags, map[string]interface{}{"action": "reject"}),
	)

	conn.Close()
	return conn, nil
}

// eagerTLSListener is a net.Listener that wraps a TLS listener and eagerly completes a TLS
// handshake with the client when a connection is accepted.
type eagerTLSListener struct {
	name             string
	handshakeTimeout time.Duration
	net.Listener
}

// newEagerTLSListener creates an eager TLS listener.
func newEagerTLSListener(name string, handshakeTimeout time.Duration, ln net.Listener) net.Listener {
	return &eagerTLSListener{
		name:             name,
		handshakeTimeout: handshakeTimeout,
		Listener:         ln,
	}
}

// Accept defers to the underlying net.Listener and additionally requests completion of a TLS
// handshake immediately. Connections that fail TLS handshakes are closed immediately.
func (l *eagerTLSListener) Accept() (net.Conn, error) {
	var cancel context.CancelFunc

	tags := map[string]interface{}{"name": l.name}

	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, &util.Error{
			Namespace:    "server",
			Message:      "failed to accept on listener",
			StackedError: err,
		}
	}

	metrics.Client.Incr(metricEagerTLSListenerAccept, tags)

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		zap.L().Warn(
			"eager TLS listener rejecting non-TLS accepted connection",
			zap.Stringer("addr", conn.RemoteAddr()),
		)
		metrics.Client.Incr(metricEagerTLSListenerHandshakeError, tags)

		conn.Close()
		return conn, nil
	}

	stopwatch := lib.NewStopwatch()
	ctx := context.Background()
	if l.handshakeTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, l.handshakeTimeout)
		defer cancel()
	}

	if err := tlsConn.HandshakeContext(ctx); err != nil {
		zap.L().Warn(
			"error during client TLS handshake",
			zap.Stringer("addr", conn.RemoteAddr()),
			zap.Duration("timeout", l.handshakeTimeout),
			zap.Error(err),
		)
		metrics.Client.Incr(metricEagerTLSListenerHandshakeError, tags)

		conn.Close()
		return conn, nil
	}

	metrics.Client.Timing(metricEagerTLSListenerHandshakeDuration, stopwatch.Elapsed(), tags)

	return conn, nil
}
