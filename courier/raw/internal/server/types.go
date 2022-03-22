package server

import (
	"courier/internal/middleware"
	"courier/internal/middleware/filters"
)

// serverContextKey is used for identifying context values in an http.Server instance.
type serverContextKey int

const (
	// ctxLocalConn references the local net.Conn connection used to serve an HTTP request.
	ctxLocalConn serverContextKey = iota
)

var (
	// filterFactories maps filter identifiers to factories that can be used to create a filter
	// instance from its config object.
	filterFactories = map[string]middleware.FilterFactory{
		"authentication":  filters.NewAuthentication,
		"body":            filters.NewBody,
		"compression":     filters.NewCompression,
		"concurrency":     filters.NewConcurrencyLimit,
		"header":          filters.NewHeader,
		"identity":        filters.NewIdentity,
		"instrumentation": filters.NewInstrumentation,
		"ip":              filters.NewIP,
		"log":             filters.NewLog,
		"metadata":        filters.NewMetadata,
		"redirect":        filters.NewRedirect,
		"rewrite":         filters.NewRewrite,
		"size":            filters.NewSizeLimit,
		"static":          filters.NewStatic,
	}
)
