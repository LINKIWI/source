package middleware

import (
	"context"
	"net/http"

	"lib.kevinlin.info/aperture/lib"

	"courier/internal/metrics"
)

const (
	metricFilterCardinality      = "filter.meta.cardinality"
	metricFilterRequestProcess   = "filter.meta.request.process"
	metricFilterRequestDuration  = "filter.meta.request.duration"
	metricFilterResponseProcess  = "filter.meta.response.process"
	metricFilterResponseDuration = "filter.meta.response.duration"
)

// composedFilter is a Filter that composes multiple underlying filters in series.
type composedFilter struct {
	filters []Filter
}

// NewComposedFilter creates a Filter that composes zero or more underlying filters.
func NewComposedFilter(filters ...Filter) Filter {
	return NewInstrumentedFilter("composed", &composedFilter{filters})
}

// ProcessRequest chains the request through each filter sequentially.
func (c *composedFilter) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	var processed []Filter
	var resp *http.Response
	var err error

	proxyReq := clientReq
	tags := map[string]interface{}{
		"filter":     "composed",
		"route_host": proxyReq.Host,
		"method":     proxyReq.Method,
		"protocol":   proxyReq.Proto,
	}

	defer func() {
		metrics.Client.Size(metricFilterCardinality, int64(len(processed)), tags)
	}()

	// Pass the request incrementally into all filters in order, but terminate early if any of
	// them errors or requests a response be sent back to the client immediately.
	for _, filter := range c.filters {
		proxyReq, resp, err = filter.ProcessRequest(proxyReq)

		if resp == nil && err == nil {
			// As long as the request processing didn't error and was not terminated
			// early (i.e. the request was either passed along without error), it should
			// be queued for future response processing.
			// Prepend the successfully processed filter so that response processing can
			// occur in reverse stack order after the reverse proxy completes.
			processed = append([]Filter{filter}, processed...)
		}

		if proxyReq != nil {
			proxyReq = proxyReq.WithContext(
				context.WithValue(
					proxyReq.Context(),
					ctxComposedProcessedFilters,
					processed,
				),
			)
		}

		if resp != nil && resp.Request != nil {
			resp.Request = resp.Request.WithContext(
				context.WithValue(
					resp.Request.Context(),
					ctxComposedProcessedFilters,
					processed,
				),
			)
		}

		if err != nil || resp != nil {
			return proxyReq, resp, err
		}
	}

	return proxyReq, resp, err
}

// ProcessResponse chains the response through each filter sequentially.
func (c *composedFilter) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	var err error

	clientResp := proxyResp

	processed, ok := clientResp.Request.Context().Value(ctxComposedProcessedFilters).([]Filter)
	if !ok {
		return clientResp, nil
	}

	for _, filter := range processed {
		clientResp, err = filter.ProcessResponse(clientResp)
		if err != nil {
			return clientResp, err
		}
	}

	return clientResp, err
}

// asyncFilter wraps a Filter to execute processing asynchronously. It is useful for logic that
// does not modify the underlying request and response and should be executed off the critical path.
type asyncFilter struct {
	Filter
}

// NewAsyncFilter creates a Filter that runs response processing asynchronously.
func NewAsyncFilter(filter Filter) Filter {
	return &asyncFilter{filter}
}

// ProcessResponse runs the underlying filter's ProcessResponse asynchronously and immediately
// returns the proxy response as-is. Note that the underlying filter must not modify the response.
func (a *asyncFilter) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	go a.Filter.ProcessResponse(proxyResp)

	return proxyResp, nil
}

// instrumentedFilter wraps a Filter with transparent instrumentation.
type instrumentedFilter struct {
	name string
	Filter
}

// NewInstrumentedFilter creates a Filter that emits metrics on the underlying filter's processing.
func NewInstrumentedFilter(name string, filter Filter) Filter {
	return &instrumentedFilter{
		name:   name,
		Filter: filter,
	}
}

// ProcessRequest runs the underlying filter's ProcessRequest and emits a counter and timer.
func (i *instrumentedFilter) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	stopwatch := lib.NewStopwatch()
	tags := map[string]interface{}{
		"filter":     i.name,
		"route_host": clientReq.Host,
		"method":     clientReq.Method,
		"protocol":   clientReq.Proto,
	}

	defer func() {
		metrics.Client.Timing(metricFilterRequestDuration, stopwatch.Elapsed(), tags)
	}()

	metrics.Client.Incr(metricFilterRequestProcess, tags)

	return i.Filter.ProcessRequest(clientReq)
}

// ProcessResponse runs the underlying filter's ProcessResponse and emits a counter and timer.
func (i *instrumentedFilter) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	stopwatch := lib.NewStopwatch()
	tags := map[string]interface{}{
		"filter":     i.name,
		"route_host": proxyResp.Request.Host,
		"method":     proxyResp.Request.Method,
		"protocol":   proxyResp.Request.Proto,
	}

	defer func() {
		metrics.Client.Timing(metricFilterResponseDuration, stopwatch.Elapsed(), tags)
	}()

	metrics.Client.Incr(metricFilterResponseProcess, tags)

	return i.Filter.ProcessResponse(proxyResp)
}
