package filters

import (
	"fmt"
	"net/http"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	headerActionSet     = "set"
	headerActionAppend  = "append"
	headerActionRemove  = "remove"
	headerActionReplace = "replace"
	headerActionRename  = "rename"
)

// headerRule describes a header modification rule for both requests and responses.
type headerRule struct {
	Action string `yaml:"action"`
	Key    string `yaml:"key"`
	Value  string `yaml:"value"`
}

// headerParams is the configuration descriptor for the header filter.
type headerParams struct {
	Match    *matchSpec   `yaml:"match"`
	Request  []headerRule `yaml:"request"`
	Response []headerRule `yaml:"response"`
}

// header is a filter that allows the injection of static headers into both the proxied request and
// client response.
type header struct {
	params *headerParams
}

// NewHeader creates a new header filter.
func NewHeader(cfg *config.Filter) (middleware.Filter, error) {
	var params headerParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse header filter params",
			StackedError: err,
		}
	}

	for _, rule := range append(params.Request, params.Response...) {
		switch rule.Action {
		case headerActionSet, headerActionAppend, headerActionRemove, headerActionReplace, headerActionRename:
		default:
			return nil, &util.Error{
				Namespace: "filters",
				Message: fmt.Sprintf(
					"header action must be one of %v",
					[]string{headerActionSet, headerActionAppend, headerActionRemove, headerActionReplace, headerActionRename},
				),
				Tags: map[string]interface{}{
					"action": rule.Action,
					"key":    rule.Key,
					"value":  rule.Value,
				},
			}
		}
	}

	filter = &header{params: &params}
	if params.Match != nil {
		filter = newMatch("identity", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("header", filter), nil
}

// ProcessRequest modifies the headers in the proxy request bound for the upstream.
func (h *header) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	for _, rule := range h.params.Request {
		switch rule.Action {
		case headerActionSet:
			clientReq.Header.Set(rule.Key, rule.Value)
		case headerActionAppend:
			clientReq.Header.Add(rule.Key, rule.Value)
		case headerActionRemove:
			clientReq.Header.Del(rule.Key)
		case headerActionReplace:
			if clientReq.Header.Get(rule.Key) != "" {
				clientReq.Header.Set(rule.Key, rule.Value)
			}
		case headerActionRename:
			if value := clientReq.Header.Get(rule.Key); value != "" {
				clientReq.Header.Set(rule.Value, value)
				clientReq.Header.Del(rule.Key)
			}
		}
	}

	return clientReq, nil, nil
}

// ProcessResponse modifies the headers in the client response bound for the client.
func (h *header) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	if proxyResp.Header == nil {
		proxyResp.Header = make(http.Header)
	}

	for _, rule := range h.params.Response {
		switch rule.Action {
		case headerActionSet:
			proxyResp.Header.Set(rule.Key, rule.Value)
		case headerActionAppend:
			proxyResp.Header.Add(rule.Key, rule.Value)
		case headerActionRemove:
			proxyResp.Header.Del(rule.Key)
		case headerActionReplace:
			if proxyResp.Header.Get(rule.Key) != "" {
				proxyResp.Header.Set(rule.Key, rule.Value)
			}
		case headerActionRename:
			if value := proxyResp.Header.Get(rule.Key); value != "" {
				proxyResp.Header.Set(rule.Value, value)
				proxyResp.Header.Del(rule.Key)
			}
		}
	}

	return proxyResp, nil
}
