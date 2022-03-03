package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"syscall"

	stdproxy "golang.org/x/net/proxy"
	"lib.kevinlin.info/aperture"

	"bitflow/internal/config"
	"bitflow/internal/log"
	"bitflow/pkg/proxy"
)

// Server abstracts multiple TCP proxy servers.
type Server struct {
	Config  config.Server
	Metrics aperture.Statsd
	Log     *log.Context

	proxies []*instance
}

// Serve starts all proxies and serves indefinitely.
func (s *Server) Serve() error {
	var proxies sync.WaitGroup

	errs := make(chan error, len(s.Config.Proxies))

	for _, p := range s.Config.Proxies {
		s.Log.Info.Printf(
			"server: starting proxy server instance: name=%s listener=%v",
			p.Name,
			p.Listener.Address,
		)

		inst, err := newProxyInstance(p, s.Metrics, s.Log)
		if err != nil {
			return err
		}

		s.proxies = append(s.proxies, inst)

		go func() {
			proxies.Add(1)
			defer proxies.Done()

			errs <- inst.serve()
		}()
	}

	if err := <-errs; err != nil {
		s.Log.Error.Printf(
			"server: proxy serving returned fatal error: err=%v",
			err,
		)

		return err
	}

	proxies.Wait()

	return nil
}

// Close closes all individual proxy server instances associated with this server.
func (s *Server) Close() {
	for _, inst := range s.proxies {
		s.Log.Info.Printf(
			"server: gracefully shutting down proxy server instance: "+
				"name=%s listener=%v",
			inst.cfg.Name,
			inst.cfg.Listener.Address,
		)

		if err := inst.close(); err != nil {
			s.Log.Warn.Printf(
				"server: error closing proxy server: name=%s err=%v",
				inst.cfg.Name,
				err,
			)
		}
	}
}

// instance describes a single proxy instance for an individual listener and target.
type instance struct {
	cfg           config.Proxy
	proxy         *proxy.Proxy
	listener      net.Listener
	connectionLog io.WriteCloser
}

// newProxyInstance initializes a new proxy server instance by starting a listener and creating a
// factory for dialing its configured target.
func newProxyInstance(
	cfg config.Proxy,
	metrics aperture.Statsd,
	logCtx *log.Context,
) (*instance, error) {
	var err error
	var connectionLog io.WriteCloser
	var tlsCfg *tls.Config

	if cfg.Listener.TLSContext != nil {
		logCtx.Debug.Printf(
			"server: enabling TLS termination for proxy: "+
				"name=%s listener=%v certificates=%v",
			cfg.Name,
			cfg.Listener.Address,
			cfg.Listener.TLSContext.Certificates,
		)

		tlsCfg, err = cfg.Listener.TLSContext.ServerConfig()
		if err != nil {
			return nil, fmt.Errorf(
				"config: error creating server TLS configuration: err=%v",
				err,
			)
		}
	}

	ln, err := cfg.Listener.Address.Listen(tlsCfg)
	if err != nil {
		return nil, fmt.Errorf("server: error listening on socket: err=%v", err)
	}

	if len(cfg.Listener.AuthorizedSources) > 0 {
		logCtx.Debug.Printf(
			"server: enabling restricted client sources for proxy: "+
				"name=%s listener=%v sources=%v",
			cfg.Name,
			cfg.Listener.Address,
			cfg.Listener.AuthorizedSources,
		)

		ln = &restrictedListener{
			sources:  cfg.Listener.AuthorizedSources,
			Listener: ln,
		}
	}

	if cfg.Options.ConnectionLog != "" {
		logCtx.Debug.Printf(
			"server: enabling connection log for proxy: name=%s listener=%v log=%s",
			cfg.Name,
			cfg.Listener.Address,
			cfg.Options.ConnectionLog,
		)

		connectionLog, err = os.OpenFile(
			cfg.Options.ConnectionLog,
			syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT,
			0644,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"server: error opening connection log file: err=%v",
				err,
			)
		}
	}

	p := &proxy.Proxy{
		Name:                 cfg.Name,
		MaxActiveConnections: cfg.Options.ConnectionLimit,
		EnableProxyHeader:    cfg.Options.EnableProxyHeader,
		DialTarget:           newMultiTargetDialer(cfg, logCtx),
		Metrics:              metrics,
		DebugLog:             logCtx.Debug,
		ErrorLog:             logCtx.Error,
		ConnectionLog:        connectionLog,
	}

	return &instance{
		cfg:           cfg,
		proxy:         p,
		listener:      ln,
		connectionLog: connectionLog,
	}, nil
}

// newMultiTargetDialer is a dialer implementation that abstracts over several candidate targets
// with a load balancing policy.
func newMultiTargetDialer(cfg config.Proxy, logCtx *log.Context) targetConnFactory {
	candidates := make([]targetCandidate, len(cfg.Targets))

	for idx, target := range cfg.Targets {
		candidates[idx] = func(target *config.Target) targetCandidate {
			factory := func(src net.Conn) (net.Conn, proxy.Pipe, error) {
				dialer := newDirectTargetDialer(target, cfg.Options.EnableDualStack)

				conn, err := dialer(src)
				if err != nil {
					return nil, nil, err
				}

				name := fmt.Sprintf("%s/%s", cfg.Name, target.Name)
				pipe := &proxy.DirectPipe{
					Name:                    name,
					ConnectionLifetime:      cfg.Options.ConnectionLifetime,
					SourceReadTimeout:       cfg.Listener.ReadTimeout,
					SourceWriteTimeout:      cfg.Listener.WriteTimeout,
					DestinationReadTimeout:  target.ReadTimeout,
					DestinationWriteTimeout: target.WriteTimeout,
					DebugLog:                logCtx.Debug,
					ErrorLog:                logCtx.Error,
				}

				return conn, pipe, nil
			}

			return targetCandidate{
				cfg:     target,
				factory: factory,
			}
		}(target)
	}

	switch cfg.Options.ConnectionLoadBalancer {
	case "failover":
		return newFailoverLoadBalancer(candidates)
	case "round-robin":
		return newRoundRobinBalancer(candidates)
	case "sni":
		return newSNILoadBalancer(candidates)
	default:
		return newStaticLoadBalancer(candidates)
	}
}

// newDirectTargetDialer is a dialer implementation factory for directly dialing the specified
// target address on-demand.
func newDirectTargetDialer(
	target *config.Target,
	enableDualStack bool,
) func(net.Conn) (net.Conn, error) {
	return func(src net.Conn) (net.Conn, error) {
		var err error
		var dialer stdproxy.Dialer
		var srcIPNet string

		ctx := context.Background()

		network, addr := target.Address.Address()
		dialer = &net.Dialer{
			Timeout:   target.ConnectTimeout,
			KeepAlive: target.KeepaliveInterval,
		}

		if target.TLSContext != nil {
			tlsCfg, err := target.TLSContext.ClientConfig()
			if err != nil {
				return nil, fmt.Errorf(
					"server: error creating client TLS configuration: err=%v",
					err,
				)
			}

			dialer = &tls.Dialer{
				NetDialer: dialer.(*net.Dialer),
				Config:    tlsCfg,
			}
		}

		if enableDualStack {
			if srcTCPAddr, ok := src.RemoteAddr().(*net.TCPAddr); ok {
				switch {
				case srcTCPAddr.IP.To4() != nil:
					srcIPNet = "ip4"
				case srcTCPAddr.IP.To16() != nil:
					srcIPNet = "ip6"
				default:
					srcIPNet = "ip"
				}
			}

			if target.ResolveTimeout > 0 {
				var cancel context.CancelFunc

				ctx, cancel = context.WithTimeout(ctx, target.ResolveTimeout)
				defer cancel()
			}

			resolved, err := target.Address.Resolve(ctx, srcIPNet)
			if err != nil {
				return nil, fmt.Errorf(
					"server: error resolving target address for dual stack "+
						"selection: err=%v",
					err,
				)
			}

			network, addr = resolved.Address()
		}

		if target.Proxy != nil {
			proxyNet, proxyAddr := target.Proxy.Address()

			if enableDualStack {
				proxyResolved, err := target.Proxy.Resolve(ctx, srcIPNet)
				if err != nil {
					return nil, fmt.Errorf(
						"server: error resolving proxy address for dual "+
							"stack selection: err=%v",
						err,
					)
				}

				proxyNet, proxyAddr = proxyResolved.Address()
			}

			dialer, err = stdproxy.SOCKS5(proxyNet, proxyAddr, nil, dialer)
			if err != nil {
				return nil, fmt.Errorf("server: error dialing proxy: err=%v", err)
			}
		}

		if target.ConnectAttempts > 0 {
			dialer = &retryDialer{
				attempts: target.ConnectAttempts,
				backoff:  target.ConnectBackoff,
				Dialer:   dialer,
			}
		}

		return dialer.Dial(network, addr)
	}
}

// serve starts serving the proxy server.
func (i *instance) serve() error {
	return i.proxy.Serve(i.listener)
}

// close requests asynchronous graceful shutdown of the proxy instance by requesting release of all
// resources held by dependencies.
func (i *instance) close() error {
	if i.connectionLog != nil {
		i.connectionLog.Close()
	}

	// Request graceful shutdown of the proxy followed by closing the associated listener.
	// After the listener is closed, Proxy's Serve will return after all connections have been
	// drained.
	i.proxy.Close()
	return i.listener.Close()
}
