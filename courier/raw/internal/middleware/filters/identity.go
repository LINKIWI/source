package filters

import (
	"net/http"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/middleware"
	"courier/internal/util"
)

// identityRule describes a single rule for a client identity pattern.
type identityRule struct {
	Action  authorizationAction `yaml:"action"`
	Pattern *config.Regex       `yaml:"pattern"`
	Serial  *config.Regex       `yaml:"serial"`
}

// identityParams is the configuration descriptor for the identity filter.
type identityParams struct {
	Match   *matchSpec          `yaml:"match"`
	Default authorizationAction `yaml:"default"`
	Rules   []identityRule      `yaml:"rules"`
}

// NewIdentity creates a new identity filter.
func NewIdentity(cfg *config.Filter) (middleware.Filter, error) {
	var params identityParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse identity filter params",
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
					"pattern":   rule.Pattern,
					"action":    params.Default,
					"available": []authorizationAction{actionAllow, actionDeny},
				},
			}
		}
	}

	filter = newAuthorization("identity", &identityAuthorizer{
		defaultAction: params.Default,
		rules:         params.Rules,
	})

	if params.Match != nil {
		filter = newMatch("identity", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("identity", filter), nil
}

// identityAuthorizer implements the authorizer interface to authorize requests by the identity
// expressed in the DNS names of supplied valid client certificates.
type identityAuthorizer struct {
	defaultAction authorizationAction
	rules         []identityRule
}

// RequestAuthorized extracts all client TLS certificates and compares all DNS names against
// patterns specified in configuration rules, executing the corresponding action if there is a
// match.
func (i *identityAuthorizer) RequestAuthorized(request *http.Request) authorizationAction {
	if request.TLS == nil || len(request.TLS.PeerCertificates) == 0 {
		zap.L().Debug(
			"identity authorization enabled but no client certificates are available",
			zap.Any("action", i.defaultAction),
		)

		return i.defaultAction
	}

	for _, rule := range i.rules {
		if rule.Serial != nil {
			serial := request.TLS.PeerCertificates[0].SerialNumber.String()

			if rule.Serial.MatchString(serial) {
				return rule.Action
			}
		}

		if rule.Pattern != nil {
			for _, name := range request.TLS.PeerCertificates[0].DNSNames {
				if rule.Pattern.MatchString(name) {
					return rule.Action
				}
			}
		}
	}

	return i.defaultAction
}
