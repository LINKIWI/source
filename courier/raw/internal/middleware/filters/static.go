package filters

import (
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

// staticParams is the configuration descriptor for the static filter.
type staticParams struct {
	Match    *matchSpec `yaml:"match"`
	Response struct {
		Status  int `yaml:"status"`
		Headers []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		} `yaml:"headers"`
		Body string `yaml:"body"`
	} `yaml:"response"`
}

// static is a filter that returns to the client a statically defined HTTP response when specified
// request parameters are matched.
type static struct {
	params *staticParams
	noop
}

// NewStatic creates a new static filter.
func NewStatic(cfg *config.Filter) (middleware.Filter, error) {
	var params staticParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse static filter params",
			StackedError: err,
		}
	} else if params.Response.Status != 0 && http.StatusText(params.Response.Status) == "" {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "invalid HTTP response status code",
			Tags:      map[string]interface{}{"status": params.Response.Status},
		}
	}

	for _, hdr := range params.Response.Headers {
		if hdr.Key == "" {
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "empty header key",
			}
		} else if hdr.Value == "" {
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "empty header value",
				Tags:      map[string]interface{}{"key": hdr.Key},
			}
		}
	}

	filter = &static{params: &params}
	if params.Match != nil {
		filter = newMatch("static", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("static", filter), nil
}

// ProcessRequest checks the request parameters against all non-nil match criteria, and terminates
// early with the predefined static response if all criteria are met. When there are no matches or
// incomplete matches, the request is returned as-is for propagation through the remaining filters.
func (s *static) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	headers := make(http.Header)
	for _, hdr := range s.params.Response.Headers {
		headers.Set(hdr.Key, hdr.Value)
	}

	status := http.StatusOK
	if s.params.Response.Status != 0 {
		status = s.params.Response.Status
	}

	resp := &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(s.params.Response.Body)),
		Request:    clientReq,
		Header:     headers,
	}

	zap.L().Debug(
		"static filter issued a static response",
		zap.Int("status", status),
	)

	return nil, resp, nil
}
