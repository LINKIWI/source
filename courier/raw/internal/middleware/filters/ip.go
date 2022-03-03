package filters

import (
	"net"
	"net/http"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

// ipRule describes a single rule for a CIDR block.
type ipRule struct {
	Action    authorizationAction `yaml:"action"`
	Transport *config.Regex       `yaml:"transport"`
	CIDR      *config.CIDR        `yaml:"cidr"`
}

// ipParams is the configuration descriptor for the ip filter.
type ipParams struct {
	Match   *matchSpec          `yaml:"match"`
	Default authorizationAction `yaml:"default"`
	Rules   []ipRule            `yaml:"rules"`
}

// NewIP creates a new IP authorization filter.
func NewIP(cfg *config.Filter) (middleware.Filter, error) {
	var params ipParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse IP filter params",
			StackedError: err,
		}
	}

	switch params.Default {
	case actionAllow, actionDeny:
	default:
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "supplied default action invalid",
			Tags: map[string]interface{}{
				"action":    params.Default,
				"available": []authorizationAction{actionAllow, actionDeny},
			},
		}
	}

	for _, rule := range params.Rules {
		switch rule.Action {
		case actionAllow, actionDeny:
		default:
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "supplied rule action invalid",
				Tags: map[string]interface{}{
					"cidr":      rule.CIDR,
					"action":    params.Default,
					"available": []authorizationAction{actionAllow, actionDeny},
				},
			}
		}
	}

	filter = newAuthorization("ip", &ipAuthorizer{
		defaultAction: params.Default,
		rules:         params.Rules,
	})

	if params.Match != nil {
		filter = newMatch("ip", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("ip", filter), nil
}

// ipAuthorizer implements the authorizer interface to authorize requests by client IP address.
type ipAuthorizer struct {
	defaultAction authorizationAction
	rules         []ipRule
}

// RequestAuthorized evaluates the remote address against each CIDR block rule and executes the
//// corresponding action if there is a match.
func (i *ipAuthorizer) RequestAuthorized(request *http.Request) authorizationAction {
	for _, rule := range i.rules {
		switch {
		case rule.Transport == nil,
			rule.Transport.MatchString("tcp"),
			rule.Transport.MatchString("tcp4"),
			rule.Transport.MatchString("tcp6"):
			clientIP, _, err := net.SplitHostPort(request.RemoteAddr)
			if err != nil {
				zap.L().Error(
					"failed to parse client remote address",
					zap.String("address", request.RemoteAddr),
					zap.Error(err),
				)
				return i.defaultAction
			}

			if rule.CIDR.Contains(net.ParseIP(clientIP)) {
				return rule.Action
			}
		case rule.Transport.MatchString("unix"):
			if request.RemoteAddr == "@" {
				return rule.Action
			}
		}
	}

	return i.defaultAction
}
