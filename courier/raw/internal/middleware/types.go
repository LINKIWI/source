package middleware

import (
	"net/http"

	"courier/internal/config"
)

// metaFilterContextKey is used for arbitrary state attached to a request, intended for use by meta
// filters.
type metaFilterContextKey int

// FilterFactory creates a filter from a config.Filter node.
type FilterFactory func(cfg *config.Filter) (Filter, error)

// Filter is an interface for logic that processes a proxied HTTP request.
type Filter interface {
	// ProcessRequest processes the inbound client request before it is proxied to the upstream
	// server. It may modify the request before it is sent or terminate early with a response.
	ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error)

	// ProcessResponse processes the proxied response from the upstream server before it is sent
	// back to the client. It may modify the response before it is sent.
	ProcessResponse(proxyResp *http.Response) (*http.Response, error)
}

const (
	ctxComposedProcessedFilters metaFilterContextKey = iota
)
