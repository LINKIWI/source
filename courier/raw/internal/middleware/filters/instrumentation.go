package filters

import (
	"context"
	"net/http"
	"strings"

	"lib.kevinlin.info/aperture/lib"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricRequestProcess   = "proxy.request.process"
	metricRequestDuration  = "proxy.request.duration"
	metricRequestBodySize  = "proxy.request.body_size"
	metricResponseBodySize = "proxy.response.body_size"
)

// instrumentationParams is the configuration descriptor for the instrumentation filter.
type instrumentationParams struct {
	Tags []struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"tags"`
}

// instrumentation is a filter that emits useful metrics on requests and responses.
type instrumentation struct {
	params *instrumentationParams
}

// NewInstrumentation creates a new instrumentation filter.
func NewInstrumentation(cfg *config.Filter) (middleware.Filter, error) {
	var params instrumentationParams

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse instrumentation filter params",
			StackedError: err,
		}
	}

	filter := &instrumentation{params: &params}

	return middleware.NewInstrumentedFilter("instrumentation", filter), nil
}

// ProcessRequest starts a latency stopwatch for the request and wraps the request body with size
// tracking.
func (i *instrumentation) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	reqBody := util.NewInterceptedIOStream(clientReq.Body)

	ctx := context.WithValue(clientReq.Context(), ctxInstrumentationStopwatch, lib.NewStopwatch())
	ctx = context.WithValue(ctx, ctxInstrumentationRequestBody, reqBody)

	clientReq.Body = reqBody

	return clientReq.WithContext(ctx), nil, nil
}

// ProcessResponse emits metrics about the request/response pair.
func (i *instrumentation) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	var peerIdentity []string

	stopwatch := proxyResp.Request.Context().Value(ctxInstrumentationStopwatch).(*lib.Stopwatch)
	reqBody := proxyResp.Request.Context().Value(ctxInstrumentationRequestBody).(*util.InterceptedIOStream)

	tags := map[string]interface{}{
		"status":     proxyResp.StatusCode,
		"route_host": proxyResp.Request.Host,
		"method":     proxyResp.Request.Method,
		"protocol":   proxyResp.Request.Proto,
		"tls":        proxyResp.Request.TLS != nil,
	}

	if proxyResp.Request.TLS != nil {
		for _, peerCert := range proxyResp.Request.TLS.PeerCertificates {
			for _, name := range peerCert.DNSNames {
				peerIdentity = append(peerIdentity, name)
			}
		}

		tags["peer"] = strings.Join(peerIdentity, ",")
	}

	for _, tag := range i.params.Tags {
		tags[tag.Key] = tag.Value
	}

	// By the time the proxy response is passed through this filter, the original request body
	// has necessarily been closed.
	if proxyResp.Request.Body != nil {
		metrics.Client.Size(metricRequestBodySize, int64(reqBody.BytesRead()), tags)
	}

	// The response body will be written back to the client at some later time, after the entire
	// filter chain has been executed. This sets a callback to report the response body size
	// after it is closed by the governing http.ResponseWriter.
	respBody := util.NewInterceptedIOStream(proxyResp.Body)
	respBody.RegisterCallback(func(responseBytes int, _ int) {
		metrics.Client.Size(metricResponseBodySize, int64(responseBytes), tags)
	})
	proxyResp.Body = respBody

	metrics.Client.Incr(metricRequestProcess, tags)
	metrics.Client.Timing(metricRequestDuration, stopwatch.Elapsed(), tags)

	return proxyResp, nil
}
