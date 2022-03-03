package filters

import (
	"net/http"

	"courier/internal/config"
	"courier/internal/middleware"
)

// internal is a default filter used for proxy-internal logic.
type internal struct {
	noop
}

// NewInternal creates a new internal filter.
func NewInternal(cfg *config.Filter) (middleware.Filter, error) {
	return &internal{}, nil
}

func (i *internal) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	// Populate inbound request metadata in the request URL

	clientReq.URL.Scheme = "http"
	if clientReq.TLS != nil {
		clientReq.URL.Scheme = "https"
	}

	clientReq.URL.Host = clientReq.Host

	return clientReq, nil, nil
}
