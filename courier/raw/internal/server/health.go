package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"lib.kevinlin.info/aperture/lib"

	"courier/internal/config"
	"courier/internal/meta"
	"courier/internal/metrics"
	"courier/internal/util"
)

const (
	metricHealthCheckTotal    = "upstream.health_check.total"
	metricHealthCheckEvaluate = "upstream.health_check.evaluate"
	metricHealthCheckDuration = "upstream.health_check.duration"
)

// healthProbe is an asynchronous ticker that periodically probes the health of an upstream server.
type healthProbe struct {
	vhost    string
	upstream *config.Upstream
	client   *http.Client
	err      error
}

// newHealthProbe creates a new health probe given an identifier and upstream configuration.
func newHealthProbe(vhost string, upstream *config.Upstream) (*healthProbe, error) {
	transport, err := upstream.ClientTransport()
	if err != nil {
		return nil, &util.Error{
			Namespace:    "server",
			Message:      "failed to create client transport",
			StackedError: err,
		}
	}

	hp := &healthProbe{
		vhost:    vhost,
		upstream: upstream,
		client: &http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	go hp.start()

	return hp, nil
}

// error reports the current error associated with the health probe.
// Semantically, a nil error indicates that the upstream is healthy; a non-nil error indicates that
// the upstream is unhealthy.
func (hp *healthProbe) error() error {
	return hp.err
}

// start starts the background periodic upstream health probes, updating the current error on every
// probe.
func (hp *healthProbe) start() {
	var attempts int

	tags := map[string]interface{}{
		"virtual_host": hp.vhost,
		"upstream":     hp.upstream.Name,
		"address":      hp.upstream.Address,
		"transport":    hp.upstream.Transport,
		"method":       hp.upstream.HealthCheck.Method,
		"path":         hp.upstream.HealthCheck.Path,
		"tls":          hp.upstream.TLSContext != nil,
	}

	for range time.NewTicker(hp.upstream.HealthCheck.Interval).C {
		attempts++

		zap.L().Debug(
			"executing health check against upstream",
			zap.String("upstream", hp.upstream.Name),
			zap.Stringer("address", hp.upstream.Address),
			zap.String("host", hp.upstream.HealthCheck.Host),
			zap.String("method", hp.upstream.HealthCheck.Method),
			zap.String("path", hp.upstream.HealthCheck.Path),
		)

		if hp.upstream.HealthCheck.Jitter > 0 {
			time.Sleep(time.Duration(float64(hp.upstream.HealthCheck.Jitter) * rand.Float64()))
		}

		stopwatch := lib.NewStopwatch()

		if err := hp.probe(); err != nil {
			hp.err = err
			metrics.Client.Incr(
				metricHealthCheckEvaluate,
				util.MergeMaps(tags, map[string]interface{}{"result": "fail"}),
			)
			zap.L().Debug(
				"upstream health check fail",
				zap.String("upstream", hp.upstream.Name),
				zap.Stringer("address", hp.upstream.Address),
				zap.String("host", hp.upstream.HealthCheck.Host),
				zap.String("method", hp.upstream.HealthCheck.Method),
				zap.String("path", hp.upstream.HealthCheck.Path),
				zap.Duration("duration", stopwatch.Elapsed()),
				zap.Error(err),
			)
		} else {
			hp.err = nil
			metrics.Client.Incr(
				metricHealthCheckEvaluate,
				util.MergeMaps(tags, map[string]interface{}{"result": "pass"}),
			)
			zap.L().Debug(
				"upstream health check pass",
				zap.String("upstream", hp.upstream.Name),
				zap.Stringer("address", hp.upstream.Address),
				zap.String("host", hp.upstream.HealthCheck.Host),
				zap.String("method", hp.upstream.HealthCheck.Method),
				zap.String("path", hp.upstream.HealthCheck.Path),
				zap.Duration("duration", stopwatch.Elapsed()),
			)
		}

		metrics.Client.Timing(metricHealthCheckDuration, stopwatch.Elapsed(), tags)
		metrics.Client.Gauge(metricHealthCheckTotal, float64(attempts), tags)
	}
}

// probe sends a single HTTP request to the upstream to check its health.
func (hp *healthProbe) probe() error {
	scheme := "http"
	host := "localhost"
	path := "/"
	if hp.upstream.TLSContext != nil {
		scheme = "https"
	}
	if hp.upstream.HealthCheck.Host != "" {
		host = hp.upstream.HealthCheck.Host
	}
	if hp.upstream.HealthCheck.Path != "" {
		path = hp.upstream.HealthCheck.Path
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	req, err := http.NewRequest(hp.upstream.HealthCheck.Method, u.String(), nil)
	if err != nil {
		return &util.Error{
			Namespace: "server",
			Message:   "probe failure due to error creating HTTP client request",
			Tags: map[string]interface{}{
				"url": u.String(),
			},
			StackedError: err,
		}
	}

	req.Header.Add(
		"User-Agent",
		fmt.Sprintf(
			"courier/%s (healthcheck; upstream:%s)",
			meta.Version,
			hp.upstream.Name,
		),
	)

	if len(hp.upstream.HealthCheck.Headers) > 0 {
		for _, header := range hp.upstream.HealthCheck.Headers {
			req.Header.Add(header.Key, header.Value)
		}
	}

	resp, err := hp.client.Do(req)
	if err != nil {
		return &util.Error{
			Namespace:    "server",
			Message:      "probe failure due to request failure",
			StackedError: err,
		}
	}

	defer resp.Body.Close()

	if hp.upstream.HealthCheck.Status != nil {
		// Custom status code range
		if resp.StatusCode < hp.upstream.HealthCheck.Status.Lower ||
			resp.StatusCode >= hp.upstream.HealthCheck.Status.Upper {
			return &util.Error{
				Namespace: "server",
				Message:   "probe failure due to invalid status not in range",
				Tags: map[string]interface{}{
					"url":    u.String(),
					"lower":  hp.upstream.HealthCheck.Status.Lower,
					"upper":  hp.upstream.HealthCheck.Status.Upper,
					"status": resp.StatusCode,
				},
			}
		}
	} else {
		// By default, assume any 2xx response to be valid
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return &util.Error{
				Namespace: "server",
				Message:   "probe failure due to invalid non-2xx status",
				Tags: map[string]interface{}{
					"url":    u.String(),
					"status": resp.StatusCode,
				},
			}
		}
	}

	return nil
}
