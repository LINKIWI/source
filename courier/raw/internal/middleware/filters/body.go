package filters

import (
	"bytes"
	"io"
	"net/http"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

// bodyRule describes a body modification rule for both requests and responses.
type bodyRule struct {
	Find    *config.Regex `yaml:"find"`
	Replace string        `yaml:"replace"`
}

// bodyParams is the configuration descriptor for the body filter.
type bodyParams struct {
	Match    *matchSpec `yaml:"match"`
	Request  []bodyRule `yaml:"request"`
	Response []bodyRule `yaml:"response"`
}

// body is a filter that performs regular expression manipulation on request and response bodies.
type body struct {
	params *bodyParams
}

// NewBody creates a new body filter.
func NewBody(cfg *config.Filter) (middleware.Filter, error) {
	var params bodyParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse body filter params",
			StackedError: err,
		}
	}

	filter = &body{params: &params}
	if params.Match != nil {
		filter = newMatch("body", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("body", filter), nil
}

// ProcessRequest buffers the entire request body in-memory, executes the request manipulation
// rules, and forwards the request to the next filter.
func (b *body) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	if len(b.params.Request) == 0 {
		return clientReq, nil, nil
	}

	buf, err := io.ReadAll(clientReq.Body)
	defer clientReq.Body.Close()

	if err != nil {
		return nil, nil, &util.Error{
			Namespace:    "filters",
			Message:      "error buffering request body",
			StackedError: err,
		}
	}

	for _, rule := range b.params.Request {
		buf = rule.Find.ReplaceAll(buf, []byte(rule.Replace))
	}

	clientReq.Body = io.NopCloser(bytes.NewReader(buf))

	return clientReq, nil, nil
}

// ProcessResponse buffers the entire response body in-memory, executes the response manipulation
// rules, and forwards the response to the next filter.
func (b *body) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	if len(b.params.Response) == 0 {
		return proxyResp, nil
	}

	// Response body manipulation is not supported for bidirectional streams, e.g. websockets or
	// HTTP/2 streams.
	if _, ok := proxyResp.Body.(io.Writer); ok {
		return proxyResp, nil
	}

	// The entire response buffer needs to be buffered in-memory so that regex replacements can
	// operate on the entire body.
	buf, err := io.ReadAll(proxyResp.Body)
	defer proxyResp.Body.Close() // The response body has already been consumed

	if err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "error buffering response body",
			StackedError: err,
		}
	}

	for _, rule := range b.params.Response {
		buf = rule.Find.ReplaceAll(buf, []byte(rule.Replace))
	}

	proxyResp.Body = io.NopCloser(bytes.NewReader(buf))

	return proxyResp, nil
}
