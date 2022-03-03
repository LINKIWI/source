package filters

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricAuthorizationEvaluate = "filter.authorization.evaluate"
)

var (
	errRequestDenied = errors.New("request denied by proxy configuration")
)

// authorizationAction describes the action to take for a particular request.
type authorizationAction string

const (
	// actionAllow indicates the request is explicitly allowed.
	actionAllow authorizationAction = "allow"
	// actionDeny indicates the request is explicitly denied.
	actionDeny authorizationAction = "deny"
)

// authorizer is an interface for expressing the authorization action to take for a client request.
type authorizer interface {
	RequestAuthorized(request *http.Request) authorizationAction
}

// authorization is a filter that authorizes individual incoming requests. Allowed requests are a
// noop; denied requests return a standard response with http.StatusForbidden.
type authorization struct {
	name string
	rule authorizer
	noop
}

// newAuthorization creates a new authorization filter with the specified authorizer implementation.
func newAuthorization(name string, rule authorizer) middleware.Filter {
	return &authorization{
		name: name,
		rule: rule,
	}
}

// ProcessRequest returns the request as-is unmodified when the authorizer indicates the request is
// allowed, or a static deny response otherwise.
func (a *authorization) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	tags := map[string]interface{}{
		"filter":     a.name,
		"route_host": clientReq.Host,
		"method":     clientReq.Method,
		"protocol":   clientReq.Proto,
	}

	denyResp := &http.Response{
		StatusCode: http.StatusForbidden,
		Body:       io.NopCloser(strings.NewReader(errRequestDenied.Error())),
		Request:    clientReq,
		Header:     make(http.Header),
	}

	action := a.rule.RequestAuthorized(clientReq)

	switch action {
	case actionAllow:
		metrics.Client.Incr(
			metricAuthorizationEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "allow"}),
		)
		return clientReq, nil, nil
	case actionDeny:
		metrics.Client.Incr(
			metricAuthorizationEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "deny"}),
		)
		return nil, denyResp, nil
	default:
		metrics.Client.Incr(
			metricAuthorizationEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "unknown"}),
		)
		zap.L().Warn(
			"authorizer returned unknown action; defaulting to deny",
			zap.String("name", a.name),
			zap.Any("action", action),
		)
		return nil, denyResp, nil
	}
}
