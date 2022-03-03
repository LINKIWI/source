package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync/atomic"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/net/netutil"

	"courier/internal/config"
	"courier/internal/log"
	"courier/internal/metrics"
	"courier/internal/util"
)

const (
	metricConnectionStateTransition = "server.connection.transition"
	metricConnectionStateTotal      = "server.connection.total"
)

// Server is the main Courier server entry point.
type Server struct {
	cfg     *config.Config
	rp      http.Handler
	servers []*instance
}

// New creates a new Courier server from the configuration.
func New(cfg *config.Config) (*Server, error) {
	rp, err := newReverseProxy(cfg.Server.VirtualHosts, cfg.Server.Proxy)
	if err != nil {
		return nil, &util.Error{
			Namespace:    "server",
			Message:      "failed to initialize reverse proxy",
			StackedError: err,
		}
	}

	return &Server{
		cfg:     cfg,
		rp:      rp,
		servers: make([]*instance, len(cfg.Server.Listeners)),
	}, nil
}

// Serve starts all individual server instances and serves indefinitely.
func (s *Server) Serve() error {
	errs := make(chan error)

	for i, listener := range s.cfg.Server.Listeners {
		s.servers[i] = newInstance(listener, s.rp)

		go func(inst *instance) {
			err := inst.serve()

			// Ignore http.ErrServerClosed errors, since they are always returned when
			// the underlying server is closed and not indicative of a real error during
			// serving. Propagate all other errors or non-errors directly as-is.
			if err == http.ErrServerClosed {
				errs <- nil
			} else {
				errs <- err
			}
		}(s.servers[i])
	}

	return <-errs
}

// Close requests graceful shutdown of all server instances.
func (s *Server) Close() error {
	for _, inst := range s.servers {
		if err := inst.close(); err != nil {
			zap.L().Warn(
				"error closing server instance",
				zap.String("name", inst.cfg.Name),
				zap.Stringer("address", inst.cfg.Address),
				zap.Error(err),
			)
		}
	}

	return nil
}

// instance represents a single server instance.
type instance struct {
	cfg    *config.Listener
	server *http.Server
}

// newInstance creates an instance wrapping a single http.Server for one listener specification.
func newInstance(cfg *config.Listener, handler http.Handler) *instance {
	connStats := util.NewConcurrentMap()
	connStats.Set(http.StateNew, new(int64))
	connStats.Set(http.StateActive, new(int64))
	connStats.Set(http.StateIdle, new(int64))
	connStats.Set(http.StateHijacked, new(int64))
	connStats.Set(http.StateClosed, new(int64))

	handler = h2c.NewHandler(handler, &http2.Server{})
	handler = log.NewRequestLogHandler(cfg.Name, handler)

	server := &http.Server{
		Handler:           handler,
		ReadTimeout:       cfg.Connection.ReadTimeout,
		ReadHeaderTimeout: cfg.Connection.ReadHeaderTimeout,
		WriteTimeout:      cfg.Connection.WriteTimeout,
		IdleTimeout:       cfg.Connection.IdleTimeout,
		ErrorLog:          log.StdErrorLogger,
		ConnState: func(conn net.Conn, state http.ConnState) {
			tags := map[string]interface{}{
				"listener": cfg.Name,
				"state":    state.String(),
			}

			stateTotal, _ := connStats.Get(state)
			atomic.AddInt64(stateTotal.(*int64), 1)

			metrics.Client.Incr(metricConnectionStateTransition, tags)
			metrics.Client.Gauge(
				metricConnectionStateTotal,
				float64(atomic.LoadInt64(stateTotal.(*int64))),
				tags,
			)
		},
		ConnContext: func(ctx context.Context, conn net.Conn) context.Context {
			return context.WithValue(ctx, ctxLocalConn, conn)
		},
	}

	return &instance{
		cfg:    cfg,
		server: server,
	}
}

// serve configures and creates the listener and starts the underlying http.Server.
func (i *instance) serve() error {
	listenNet, listenAddr := i.cfg.Address.Address()
	zap.L().Debug(
		"created server listener",
		zap.String("name", i.cfg.Name),
		zap.String("net", listenNet),
		zap.String("address", listenAddr),
	)

	ln, err := net.Listen(listenNet, listenAddr)
	if err != nil {
		return &util.Error{
			Namespace: "server",
			Message:   "failed to listen on address",
			Tags: map[string]interface{}{
				"name":    i.cfg.Name,
				"address": i.cfg.Address.String(),
			},
			StackedError: err,
		}
	}

	if i.cfg.Connection.ActiveLimit > 0 {
		zap.L().Debug(
			"enabling active connections limit for listener",
			zap.String("name", i.cfg.Name),
			zap.Int("limit", i.cfg.Connection.ActiveLimit),
		)

		ln = netutil.LimitListener(ln, i.cfg.Connection.ActiveLimit)
	}

	if len(i.cfg.AuthorizedSources) > 0 {
		for _, source := range i.cfg.AuthorizedSources {
			zap.L().Debug(
				"enabling restricted IP source for listener",
				zap.String("name", i.cfg.Name),
				zap.Stringer("source", source),
			)
		}

		ln = newRestrictedListener(i.cfg.Name, i.cfg.AuthorizedSources, ln)
	}

	switch i.cfg.Protocol {
	case "haproxy":
		zap.L().Debug(
			"enabling proxy protocol header parser for listener",
			zap.String("name", i.cfg.Name),
		)

		headerTimeout := i.cfg.Connection.ReadTimeout
		if i.cfg.Connection.ReadHeaderTimeout > 0 {
			headerTimeout = i.cfg.Connection.ReadHeaderTimeout
		}

		ln = newProxyListener(i.cfg.Name, headerTimeout, ln)
	}

	if i.cfg.TLSContext != nil {
		zap.L().Debug(
			"enabling server TLS termination",
			zap.String("name", i.cfg.Name),
			zap.Stringer("address", i.cfg.Address),
			zap.Any("context", i.cfg.TLSContext),
		)

		i.server.TLSConfig, err = i.cfg.TLSContext.ServerConfig()
		if err != nil {
			return &util.Error{
				Namespace: "server",
				Message:   "failed to build server TLS context",
				Tags: map[string]interface{}{
					"name": i.cfg.Name,
					"ctx":  i.cfg.TLSContext,
				},
				StackedError: err,
			}
		}

		ln = tls.NewListener(ln, i.server.TLSConfig)

		if i.cfg.TLSContext.HandshakeTimeout > 0 {
			ln = newEagerTLSListener(i.cfg.Name, i.cfg.TLSContext.HandshakeTimeout, ln)
		}
	}

	if err := http2.ConfigureServer(i.server, &http2.Server{}); err != nil {
		return &util.Error{
			Namespace: "server",
			Message:   "error configuring HTTP/2 on server",
			Tags: map[string]interface{}{
				"name":    i.cfg.Name,
				"address": i.cfg.Address.String(),
			},
			StackedError: err,
		}
	}

	return i.server.Serve(ln)
}

// close requests shutdown of the underlying http.Server.
func (i *instance) close() error {
	return i.server.Close()
}
