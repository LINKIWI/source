package filters

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricConcurrencyLimitEvaluate = "filter.concurrency.evaluate"
	metricConcurrencyLimitActive   = "filter.concurrency.active"
)

var (
	errConcurrencyLimitExceeded = errors.New("request concurrency limit exceeded")
)

// concurrencyLimitParams is the configuration descriptor for the concurrencyLimit filter.
type concurrencyLimitParams struct {
	MaxConcurrentRequests int64 `yaml:"max_concurrent_requests"`
}

// concurrencyLimit is a filter that limits the number of allowed concurrently serviced requests.
type concurrencyLimit struct {
	inFlightRequests *int64
	params           *concurrencyLimitParams
}

// NewConcurrencyLimit creates a new request concurrency limit filter.
func NewConcurrencyLimit(cfg *config.Filter) (middleware.Filter, error) {
	var params concurrencyLimitParams

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse request concurrency filter params",
			StackedError: err,
		}
	}

	filter := &concurrencyLimit{
		inFlightRequests: new(int64),
		params:           &params,
	}

	return middleware.NewAsyncFilter(middleware.NewInstrumentedFilter("concurrency", filter)), nil
}

// ProcessRequest rejects the inbound request if the current number of in-flight requests exceeds
// the maximum allowed, and allows the request to pass otherwise.
func (c *concurrencyLimit) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	tags := map[string]interface{}{
		"route_host": clientReq.Host,
		"method":     clientReq.Method,
		"protocol":   clientReq.Proto,
	}

	if c.params.MaxConcurrentRequests <= 0 {
		metrics.Client.Incr(
			metricConcurrencyLimitEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "noop"}),
		)
		return clientReq, nil, nil
	}

	if atomic.LoadInt64(c.inFlightRequests) >= c.params.MaxConcurrentRequests {
		metrics.Client.Incr(
			metricConcurrencyLimitEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "reject"}),
		)
		zap.L().Warn(
			"concurrency limit filter rejecting request due to quota overflow",
			zap.Int64("in_flight_current", atomic.LoadInt64(c.inFlightRequests)),
			zap.Int64("in_flight_maximum", c.params.MaxConcurrentRequests),
		)

		headers := make(http.Header)
		headers.Set(
			"X-Courier-Request-Concurrency-Limit",
			strconv.Itoa(int(c.params.MaxConcurrentRequests)),
		)

		resp := &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(strings.NewReader(errConcurrencyLimitExceeded.Error())),
			Request:    clientReq,
			Header:     headers,
		}

		return nil, resp, nil
	}

	metrics.Client.Incr(
		metricConcurrencyLimitEvaluate,
		util.MergeMaps(tags, map[string]interface{}{"action": "pass"}),
	)

	atomic.AddInt64(c.inFlightRequests, 1)

	metrics.Client.Gauge(
		metricConcurrencyLimitActive,
		float64(atomic.LoadInt64(c.inFlightRequests)),
		tags,
	)

	return clientReq, nil, nil
}

// ProcessResponse releases the in-flight request counter for the current request.
func (c *concurrencyLimit) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	tags := map[string]interface{}{
		"route_host": proxyResp.Request.Host,
		"method":     proxyResp.Request.Method,
		"protocol":   proxyResp.Request.Proto,
	}

	atomic.AddInt64(c.inFlightRequests, -1)

	metrics.Client.Gauge(
		metricConcurrencyLimitActive,
		float64(atomic.LoadInt64(c.inFlightRequests)),
		tags,
	)

	return proxyResp, nil
}
