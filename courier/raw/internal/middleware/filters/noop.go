package filters

import (
	"net/http"

	"courier/internal/config"
	"courier/internal/middleware"
)

// noop applies no processing to requests and responses.
type noop struct{}

// NewNoop creates a noop filter.
func NewNoop(_ *config.Filter) (middleware.Filter, error) {
	return &noop{}, nil
}

// ProcessRequest returns the request as-is.
func (n *noop) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	return clientReq, nil, nil
}

// ProcessResponse returns the response as-is.
func (n *noop) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	return proxyResp, nil
}
