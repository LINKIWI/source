package filters

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

// redirectParams is the configuration descriptor for the redirect filter.
type redirectParams struct {
	Match   *matchSpec    `yaml:"match"`
	Find    *config.Regex `yaml:"find"`
	Replace string        `yaml:"replace"`
	Status  int           `yaml:"status"`
}

// redirect is a filter that returns an HTTP redirect when the inbound URL matches a rule.
type redirect struct {
	params *redirectParams
	noop
}

// NewRedirect creates a new redirect filter.
func NewRedirect(cfg *config.Filter) (middleware.Filter, error) {
	var params redirectParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse redirect filter params",
			StackedError: err,
		}
	} else if params.Status != 0 && http.StatusText(params.Status) == "" {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "invalid HTTP redirect status code",
			Tags:      map[string]interface{}{"status": params.Status},
		}
	}

	filter = &redirect{params: &params}
	if params.Match != nil {
		filter = newMatch("redirect", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("redirect", filter), nil
}

// ProcessRequest attempts to match the full client URL against the redirect pattern, and if there
// is a match, executes the replacement rule on the URL and returns it as the target redirect
// location in an early termination HTTP response with an optional custom status code.
func (r *redirect) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	fullClientURL := clientReq.URL.String()

	if !r.params.Find.MatchString(fullClientURL) {
		return clientReq, nil, nil
	}

	location := r.params.Find.ReplaceAllString(fullClientURL, r.params.Replace)
	body := fmt.Sprintf("redirect: %s", location)
	headers := make(http.Header)
	headers.Set("Location", location)

	status := http.StatusTemporaryRedirect
	if r.params.Status != 0 {
		status = r.params.Status
	}

	redirectResp := &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    clientReq,
		Header:     headers,
	}

	zap.L().Debug(
		"redirect filter issued a redirect",
		zap.String("source", fullClientURL),
		zap.String("target", location),
		zap.Int("status", status),
	)

	return clientReq, redirectResp, nil
}
