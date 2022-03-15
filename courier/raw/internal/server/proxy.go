package server

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"sync/atomic"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/log"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/middleware/filters"
	"courier/internal/util"
)

const (
	metricProxyServe               = "proxy.serve"
	metricProxySuccess             = "proxy.success"
	metricProxyError               = "proxy.error"
	metricProxyTotalServed         = "proxy.serve.total"
	metricProxyEarlyTerminate      = "proxy.lifecycle.early_terminate"
	metricProxyUnknownHost         = "proxy.lifecycle.unknown_host"
	metricProxyUnknownRoute        = "proxy.lifecycle.unknown_route"
	metricProxyCircuitBreak        = "proxy.lifecycle.health_circuit_break"
	metricProxyTransportInitialize = "proxy.transport.initialize"
	metricProxyTransportRecycle    = "proxy.transport.recycle"
)

var (
	// errUnknownHost is returned when the request targets a virtual host that is not registered
	// in configuration.
	errUnknownHost = &response{
		status:  http.StatusNotFound,
		code:    "UNKNOWN_VIRTUAL_HOST",
		message: "unknown virtual host",
	}
	// errUnknownRoute is returned when the request does not match any configured route.
	errUnknownRoute = &response{
		status:  http.StatusBadRequest,
		code:    "UNKNOWN_ROUTE",
		message: "no upstream for route",
	}
	// errBadTransport is returned when the proxy encounters an error initializing the HTTP
	// client transport to an upstream.
	errBadTransport = &response{
		status:  http.StatusInternalServerError,
		code:    "BAD_TRANSPORT",
		message: "failed to initialize upstream transport",
	}
	// errFilterProcess is returned when the filter chain raises an error.
	errFilterProcess = &response{
		status:  http.StatusInternalServerError,
		code:    "FILTER_PROCESS",
		message: "failed in filter chain processing",
	}
	// errUpstreamUnhealthy is returned when the upstream health check is in a failing state, as
	// a fail-fast behavior when health check-based circuit breaking is enabled.
	errUpstreamUnhealthy = &response{
		status:  http.StatusBadGateway,
		code:    "UPSTREAM_UNHEALTHY",
		message: "upstream is marked unhealthy",
	}
	// errUpstreamTimeout is returned when the upstream service duration exceeds a provisioned
	// network-level timeout.
	errUpstreamTimeout = &response{
		status:  http.StatusGatewayTimeout,
		code:    "UPSTREAM_TIMEOUT",
		message: "upstream request timeout",
	}
	// errUpstreamUnavailable is returned for non-transient upstream network connection errors.
	errUpstreamUnavailable = &response{
		status:  http.StatusBadGateway,
		code:    "UPSTREAM_UNAVAILABLE",
		message: "upstream connection error",
	}
	// errProxyInternal is returned when the HTTP proxy itself raises an error.
	errProxyInternal = &response{
		status:  http.StatusServiceUnavailable,
		code:    "INTERNAL",
		message: "unknown internal proxy failure",
	}
	// errObscuredInternal is returned when semantically meaningful errors are disabled and the
	// proxy encounters any processing error.
	errObscuredInternal = &response{
		status:  http.StatusInternalServerError,
		code:    "OBSCURED",
		message: "internal server error",
	}
)

// internalErrorHandler is a convenience alias for httputil.ReverseProxy's ErrorHandler function.
type internalErrorHandler func(proxyRW http.ResponseWriter, proxyReq *http.Request, err error)

// reverseProxy is an http.Handler that abstracts an httputil.ReverseProxy.
type reverseProxy struct {
	vhosts []*config.VirtualHost
	opts   config.Proxy

	filterChains map[string]middleware.Filter
	healthProbes map[*config.Upstream]*healthProbe
	transports   sync.Map
	sequenceID   *int64
}

// newReverseProxy creates a new reverseProxy from a list of virtual host configurations.
func newReverseProxy(vhosts []*config.VirtualHost, opts config.Proxy) (http.Handler, error) {
	filterChains := make(map[string]middleware.Filter)
	healthProbes := make(map[*config.Upstream]*healthProbe)

	// Set up the filter chains for each virtual host up front, so that they can be trivially
	// accessed in hot request paths.
	for _, vhost := range vhosts {
		var vhostFilters []middleware.Filter

		zap.L().Debug(
			"assembling virtual host filter chain",
			zap.String("vhost", vhost.Name),
			zap.Int("filters", len(vhost.Filters)),
		)

		internal, err := filters.NewInternal(nil)
		if err != nil {
			return nil, &util.Error{
				Namespace: "server",
				Message:   "error creating internal filter",
				Tags: map[string]interface{}{
					"vhost": vhost.Name,
				},
			}
		}

		vhostFilters = append(vhostFilters, internal)

		for _, cfgFilter := range vhost.Filters {
			zap.L().Debug(
				"initializing filter",
				zap.String("vhost", vhost.Name),
				zap.String("name", cfgFilter.Name),
				zap.String("type", cfgFilter.Type),
			)

			filterFactory, ok := filterFactories[cfgFilter.Type]
			if !ok {
				return nil, &util.Error{
					Namespace: "server",
					Message:   "unknown filter type",
					Tags: map[string]interface{}{
						"vhost":  vhost.Name,
						"name":   cfgFilter.Name,
						"type":   cfgFilter.Type,
						"params": cfgFilter.Params.Value,
					},
				}
			}

			filter, err := filterFactory(cfgFilter)
			if err != nil {
				return nil, &util.Error{
					Namespace: "server",
					Message:   "error creating filter",
					Tags: map[string]interface{}{
						"vhost":  vhost.Name,
						"name":   cfgFilter.Name,
						"params": cfgFilter.Params.Value,
					},
					StackedError: err,
				}
			}

			vhostFilters = append(vhostFilters, filter)
		}

		filterChains[vhost.Name] = middleware.NewComposedFilter(vhostFilters...)

		for _, upstream := range vhost.Upstreams {
			if upstream.HealthCheck != nil {
				zap.L().Debug(
					"creating health check for virtual host upstream",
					zap.String("vhost", vhost.Name),
					zap.String("upstream", upstream.Name),
					zap.Stringer("address", upstream.Address),
					zap.String("host", upstream.HealthCheck.Host),
					zap.String("path", upstream.HealthCheck.Path),
					zap.Duration("interval", upstream.HealthCheck.Interval),
					zap.Duration("jitter", upstream.HealthCheck.Jitter),
					zap.Bool("circuit_break", upstream.HealthCheck.EnableCircuitBreak),
				)

				hp, err := newHealthProbe(vhost.Name, upstream)
				if err != nil {
					return nil, &util.Error{
						Namespace: "server",
						Message:   "error creating health check probe",
						Tags: map[string]interface{}{
							"vhost": vhost.Name,
						},
					}
				}

				healthProbes[upstream] = hp
			}
		}
	}

	return &reverseProxy{
		vhosts:       vhosts,
		opts:         opts,
		healthProbes: healthProbes,
		filterChains: filterChains,
		sequenceID:   new(int64),
	}, nil
}

// ServeHTTP serves a single HTTP request.
func (rp *reverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var err error
	var abortResp *http.Response
	var transport http.RoundTripper

	proxyReq := req.Clone(req.Context())
	tags := map[string]interface{}{
		"route_host": proxyReq.Host,
		"method":     proxyReq.Method,
		"protocol":   proxyReq.Proto,
		"tls":        proxyReq.TLS != nil,
	}

	// Match virtual host
	vhost, ok := rp.matchHost(proxyReq.Host)
	if !ok {
		metrics.Client.Incr(metricProxyUnknownHost, tags)
		zap.L().Debug("unknown virtual host", zap.String("host", proxyReq.Host))
		rp.respond(rw, errUnknownHost.response())
		return
	}

	// Match upstream from route
	upstream, ok := rp.matchRoute(proxyReq, vhost.Upstreams, vhost.Routes)
	if !ok {
		metrics.Client.Incr(metricProxyUnknownRoute, tags)
		zap.L().Debug("unknown route", zap.String("virtual_host", vhost.Name))
		rp.respond(rw, errUnknownRoute.response())
		return
	}

	tags = util.MergeMaps(tags, map[string]interface{}{
		"virtual_host": vhost.Name,
		"upstream":     upstream.Name,
	})

	metrics.Client.Incr(metricProxyServe, tags)
	metrics.Client.Gauge(metricProxyTotalServed, float64(atomic.LoadInt64(rp.sequenceID)), tags)
	defer func() { atomic.AddInt64(rp.sequenceID, 1) }()

	// Reverse proxy callbacks
	rpModifyResponse := func(proxyResp *http.Response) error {
		metrics.Client.Incr(metricProxySuccess, tags)

		proxyResp, err = rp.filterChains[vhost.Name].ProcessResponse(proxyResp)
		return err
	}
	rpErrorHandler := func(errResp *response) internalErrorHandler {
		return func(proxyRW http.ResponseWriter, proxyReq *http.Request, err error) {
			metrics.Client.Incr(
				metricProxyError,
				util.MergeMaps(tags, map[string]interface{}{"code": errResp.code}),
			)

			hub := sentry.CurrentHub().Clone()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("virtual_host", vhost.Name)
				scope.SetTag("upstream", upstream.Address.String())
				scope.SetTag("route_host", proxyReq.Host)
				scope.SetTag("method", proxyReq.Method)
				scope.SetTag("protocol", proxyReq.Proto)
				scope.SetExtra("url", proxyReq.RequestURI)
				scope.SetExtra("client", proxyReq.RemoteAddr)
				scope.SetExtra("code", errResp.code)
			})
			go hub.CaptureException(err)

			zap.L().Error(
				"reverse proxy error",
				zap.String("virtual_host", vhost.Name),
				zap.Stringer("upstream", upstream.Address),
				zap.String("route_host", proxyReq.Host),
				zap.String("method", proxyReq.Method),
				zap.String("protocol", proxyReq.Proto),
				zap.String("url", proxyReq.RequestURI),
				zap.Int("status", errResp.status),
				zap.String("code", errResp.code),
				zap.String("message", errResp.message),
				zap.Error(err),
			)

			if !rp.opts.EnableErrorSemantics {
				errResp = errObscuredInternal
			}

			if rp.opts.EnableErrorPropagation {
				errResp.tags = map[string]string{"error": err.Error()}
			}

			resp := errResp.response()
			resp.Request = proxyReq

			rpModifyResponse(resp)
			rp.respond(proxyRW, resp)
		}
	}

	// Apply request filters
	proxyReq, abortResp, err = rp.filterChains[vhost.Name].ProcessRequest(proxyReq)
	if err != nil {
		rpErrorHandler(errFilterProcess)(rw, proxyReq, err)
		return
	} else if abortResp != nil {
		metrics.Client.Incr(metricProxyEarlyTerminate, tags)
		zap.L().Debug(
			"early termination of request lifecycle before proxy",
			zap.Int("status", abortResp.StatusCode),
		)

		rpModifyResponse(abortResp)
		rp.respond(rw, abortResp)
		return
	}

	// Conditionally fail-fast on upstream health check failures
	if hp, ok := rp.healthProbes[upstream]; ok {
		if err := hp.error(); err != nil && upstream.HealthCheck.EnableCircuitBreak {
			metrics.Client.Incr(metricProxyCircuitBreak, tags)
			rpErrorHandler(errUpstreamUnhealthy)(rw, proxyReq, err)
			return
		}
	}

	// Prepare proxy request to upstream
	proxyReq.URL.Scheme = "http"
	if upstream.TLSContext != nil {
		proxyReq.URL.Scheme = "https"
	}

	// Transport initialization
	tr, ok := rp.transports.Load(upstream)
	if ok {
		transport = tr.(http.RoundTripper)
		metrics.Client.Incr(metricProxyTransportRecycle, tags)
	} else {
		transport, err = upstream.ClientTransport()
		if err != nil {
			rpErrorHandler(errBadTransport)(rw, proxyReq, err)
			return
		}

		rp.transports.Store(upstream, transport)
		metrics.Client.Incr(metricProxyTransportInitialize, tags)
	}

	backend := httputil.ReverseProxy{
		Transport: transport,
		// Instruct the proxy to flush responses back to the client immediately, for
		// compatibility with gRPC streams.
		FlushInterval: -1,
		Director: func(clientReq *http.Request) {
			// The Director routine is a noop because the request passed to the reverse
			// proxy has already been cascaded through all request filters.
		},
		ModifyResponse: rpModifyResponse,
		ErrorHandler:   rp.errorClassifierHandler(rpErrorHandler),
		ErrorLog:       log.StdWarnLogger,
	}

	backend.ServeHTTP(rw, proxyReq)
}

// matchHost attempts to match the hostname requested by the client in its Host header with a
// virtual host entry in configuration.
func (rp *reverseProxy) matchHost(reqHost string) (*config.VirtualHost, bool) {
	for _, vhost := range rp.vhosts {
		if vhost.Host.MatchString(reqHost) {
			return vhost, true
		}
	}

	return nil, false
}

// matchRoute attempts to match the inbound request to an upstream based on routing policy rules in
// configuration.
func (rp *reverseProxy) matchRoute(
	req *http.Request,
	upstreams []*config.Upstream,
	routes []*config.Route,
) (*config.Upstream, bool) {
	for _, route := range routes {
		if route.Match.Host != nil && !route.Match.Host.MatchString(req.Host) {
			continue
		} else if route.Match.Method != nil && !route.Match.Method.MatchString(req.Method) {
			continue
		} else if route.Match.Path != nil && !route.Match.Path.MatchString(req.URL.Path) {
			continue
		} else if route.Match.Header != nil {
			value := req.Header.Get(route.Match.Header.Key)
			if value == "" || !route.Match.Header.Value.MatchString(value) {
				continue
			}
		}

		for _, upstream := range upstreams {
			if upstream.Name == route.Upstream {
				return upstream, true
			}
		}
	}

	return nil, false
}

// respond writes an HTTP response from a structured response struct.
func (rp *reverseProxy) respond(rw http.ResponseWriter, resp *http.Response) error {
	if resp.Header != nil {
		for key, value := range resp.Header {
			rw.Header()[key] = value
		}
	}

	rw.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(rw, resp.Body); err != nil {
		return err
	}

	return resp.Body.Close()
}

// errorClassifierHandler is an internalErrorHandler that attempts to select the most semantically
// meaningful response for an internal error raised by the reverse proxy.
func (rp *reverseProxy) errorClassifierHandler(handler func(errResp *response) internalErrorHandler) internalErrorHandler {
	var classifier func(error) *response

	classifier = func(err error) *response {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return errUpstreamTimeout
		} else if _, ok := err.(*net.OpError); ok {
			return errUpstreamUnavailable
		} else if unwrapped := errors.Unwrap(err); unwrapped != nil {
			// Identify the first meaningful error in the error stack.
			return classifier(unwrapped)
		}

		return errProxyInternal
	}

	return func(proxyRW http.ResponseWriter, proxyReq *http.Request, err error) {
		handler(classifier(err))(proxyRW, proxyReq, err)
	}
}
