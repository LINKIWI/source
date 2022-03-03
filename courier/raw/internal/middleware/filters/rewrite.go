package filters

import (
	"net/http"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

// rewriteRule is configuration that describes how part of a URL should be rewritten.
type rewriteRule struct {
	Find    *config.Regex `yaml:"find"`
	Replace string        `yaml:"replace"`
}

// rewriteParams is the configuration descriptor for the rewrite filter.
type rewriteParams struct {
	Host []rewriteRule `yaml:"host"`
	Path []rewriteRule `yaml:"path"`
}

// rewrite is a filter that allows rewriting portions of the inbound client request before it is
// dispatched through the proxy.
type rewrite struct {
	params *rewriteParams
	noop
}

// NewRewrite creates a new rewrite filter.
func NewRewrite(cfg *config.Filter) (middleware.Filter, error) {
	var params rewriteParams

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse rewrite filter params",
			StackedError: err,
		}
	}

	return middleware.NewInstrumentedFilter("rewrite", &rewrite{params: &params}), nil
}

// ProcessRequest executes the rewrite rules against all requested components of the request URL.
// Rewrite rules are executed in the order they appear in config.
func (r *rewrite) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	if len(r.params.Host) > 0 {
		host := clientReq.Host
		for _, hostRule := range r.params.Host {
			host = hostRule.Find.ReplaceAllString(host, hostRule.Replace)
		}
		clientReq.Host = host
		clientReq.URL.Host = host
	}

	if len(r.params.Path) > 0 {
		path := clientReq.URL.Path
		for _, pathRule := range r.params.Path {
			path = pathRule.Find.ReplaceAllString(path, pathRule.Replace)
		}
		clientReq.URL.Path = path
	}

	return clientReq, nil, nil
}
