package filters

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricSizeLimitEvaluate = "filter.size.evaluate"
)

var (
	errRequestSizeLimitExceeded  = errors.New("request body size limit exceeded")
	errResponseSizeLimitExceeded = errors.New("response body size limit exceeded")
)

// sizeLimitParams is the configuration descriptor for the sizeLimit filter.
type sizeLimitParams struct {
	Request *struct {
		MaxBodySize int `yaml:"max_body_size"`
	} `yaml:"request"`
	Response *struct {
		MaxBodySize int `yaml:"max_body_size"`
	} `yaml:"response"`
}

// sizeLimit is a filter that disallows inbound requests or outbound responses beyond a configured
// maximum size threshold.
type sizeLimit struct {
	params *sizeLimitParams
}

// NewSizeLimit creates a new size limit filter.
func NewSizeLimit(cfg *config.Filter) (middleware.Filter, error) {
	var params sizeLimitParams

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse request size filter params",
			StackedError: err,
		}
	}

	if params.Request != nil && params.Request.MaxBodySize <= 0 {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "max request body size must be more than 0 bytes",
			Tags: map[string]interface{}{
				"size": params.Request.MaxBodySize,
			},
		}
	} else if params.Response != nil && params.Response.MaxBodySize <= 0 {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "max response body size must be more than 0 bytes",
			Tags: map[string]interface{}{
				"size": params.Response.MaxBodySize,
			},
		}
	}

	filter := &sizeLimit{params: &params}

	return middleware.NewInstrumentedFilter("size", filter), nil
}

// ProcessRequest rejects the inbound request if its body size exceeds the maximum allowed request
// body size. Note that the implementation buffers up to the maximum request size in-memory before
// dispatching the request to the next layer in the chain.
func (s *sizeLimit) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	// Body size limit evaluation is not supported for bidirectional streams.
	if _, ok := clientReq.Body.(io.Writer); ok {
		return clientReq, nil, nil
	}

	if s.params.Request == nil {
		return clientReq, nil, nil
	}

	tags := map[string]interface{}{
		"route_host": clientReq.Host,
		"method":     clientReq.Method,
		"protocol":   clientReq.Proto,
		"stage":      "request",
	}

	chunk := make([]byte, 1024)
	buf := make([]byte, 0, s.params.Request.MaxBodySize)

	for len(buf) < s.params.Request.MaxBodySize+1 {
		n, err := clientReq.Body.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}
	}

	if len(buf) > s.params.Request.MaxBodySize {
		metrics.Client.Incr(
			metricSizeLimitEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "reject"}),
		)
		zap.L().Warn(
			"size limit filter rejecting request due to size limit violation",
			zap.Int("max_request_size", s.params.Request.MaxBodySize),
		)

		headers := make(http.Header)
		headers.Set("X-Courier-Request-Size-Limit", strconv.Itoa(s.params.Request.MaxBodySize))

		resp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(errRequestSizeLimitExceeded.Error())),
			Request:    clientReq,
			Header:     headers,
		}

		return nil, resp, nil
	}

	metrics.Client.Incr(
		metricSizeLimitEvaluate,
		util.MergeMaps(tags, map[string]interface{}{"action": "pass"}),
	)

	// Pass along the buffered request body contents along with the remaining unread body.
	clientReq.Body = util.NewMultiReadCloser(io.NopCloser(bytes.NewReader(buf)), clientReq.Body)

	return clientReq, nil, nil
}

// ProcessResponse rejects the response if its body size exceeds the maximum allowed response body
// size. In the event that the response is rejected, its body its closed immediately and the
// client response is replaced with a static error message. Note that the implementation buffers up
// to the maximum response size in-memory before returning to the client.
func (s *sizeLimit) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	// Body size limit evaluation is not supported for bidirectional streams.
	if _, ok := proxyResp.Body.(io.Writer); ok {
		return proxyResp, nil
	}

	if s.params.Response == nil {
		return proxyResp, nil
	}

	tags := map[string]interface{}{
		"route_host": proxyResp.Request.Host,
		"method":     proxyResp.Request.Method,
		"protocol":   proxyResp.Request.Proto,
		"stage":      "response",
	}

	chunk := make([]byte, 1024)
	buf := make([]byte, 0, s.params.Response.MaxBodySize)

	for len(buf) < s.params.Response.MaxBodySize+1 {
		n, err := proxyResp.Body.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	if len(buf) > s.params.Response.MaxBodySize {
		metrics.Client.Incr(
			metricSizeLimitEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "reject"}),
		)
		zap.L().Warn(
			"size limit filter rejecting response due to size limit violation",
			zap.Int("max_response_size", s.params.Response.MaxBodySize),
		)

		// Invalidate the original response body
		proxyResp.Header.Del("Content-Length")
		proxyResp.Body.Close()

		// Replace the response with an error
		proxyResp.StatusCode = http.StatusBadRequest
		proxyResp.Header = make(http.Header)
		proxyResp.Header.Set("X-Courier-Response-Size-Limit", strconv.Itoa(s.params.Response.MaxBodySize))
		proxyResp.Body = io.NopCloser(strings.NewReader(errResponseSizeLimitExceeded.Error()))

		return proxyResp, nil
	}

	metrics.Client.Incr(
		metricSizeLimitEvaluate,
		util.MergeMaps(tags, map[string]interface{}{"action": "pass"}),
	)

	// Pass along the buffered response body contents along with the remaining unread body.
	proxyResp.Body = util.NewMultiReadCloser(io.NopCloser(bytes.NewReader(buf)), proxyResp.Body)

	return proxyResp, nil
}
