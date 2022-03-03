package filters

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	schemeBasic = "basic"
)

const (
	hashPlain  = "plain"
	hashBcrypt = "bcrypt"
	hashSHA256 = "sha256"
)

const (
	metricAuthenticationEvaluate = "filter.authentication.evaluate"
)

var (
	errAuthenticationRequired = errors.New("authentication required")
)

// authenticationParams is the configuration descriptor for the authentication filter.
type authenticationParams struct {
	Scheme     string     `yaml:"scheme"`
	Realm      string     `yaml:"realm"`
	Match      *matchSpec `yaml:"match"`
	Identities []struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Hash     string `yaml:"hash"`
	} `yaml:"identities"`
	ForwardIdentity bool `yaml:"forward_identity"`
}

// authentication is a filter supporting pluggable user authentication backends.
type authentication struct {
	params *authenticationParams
	noop
}

// NewAuthentication creates a new authentication filter.
func NewAuthentication(cfg *config.Filter) (middleware.Filter, error) {
	var params authenticationParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse authentication filter params",
			StackedError: err,
		}
	}

	switch params.Scheme {
	case schemeBasic:
	case "":
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "missing authentication scheme",
		}
	default:
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "unsupported authentication scheme",
			Tags:      map[string]interface{}{"scheme": params.Scheme},
		}
	}

	for _, id := range params.Identities {
		if id.Username == "" {
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "authentication identity username is missing",
			}
		} else if id.Password == "" {
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "authentication identity password is missing",
				Tags:      map[string]interface{}{"username": id.Username},
			}
		}

		switch id.Hash {
		case hashPlain, hashBcrypt, hashSHA256:
		case "":
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "authentication identity password hash algorithm is missing",
				Tags:      map[string]interface{}{"username": id.Username},
			}
		default:
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "unknown authentication identity password hash algorithm",
				Tags: map[string]interface{}{
					"username": id.Username,
					"hash":     id.Hash,
				},
			}
		}
	}

	filter = &authentication{params: &params}
	if params.Match != nil {
		filter = newMatch("authentication", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("authentication", filter), nil
}

// ProcessRequest selects an authentication backend based on the configuration-specified desired
// authentication scheme for this request.
func (a *authentication) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	switch a.params.Scheme {
	case schemeBasic:
		return a.processBasicAuthenticationRequest(clientReq)
	default:
		return nil, nil, &util.Error{
			Namespace: "filters",
			Message:   "unimplemented authentication scheme backend",
			Tags:      map[string]interface{}{"scheme": a.params.Scheme},
		}
	}
}

// processBasicAuthenticationRequest attempts to authenticate using HTTP basic authentication.
func (a *authentication) processBasicAuthenticationRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	tags := map[string]interface{}{
		"route_host": clientReq.Host,
		"realm":      a.params.Realm,
		"scheme":     schemeBasic,
	}

	headers := make(http.Header)
	headers.Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", a.params.Realm))
	unauthResp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(strings.NewReader(errAuthenticationRequired.Error())),
		Request:    clientReq,
		Header:     headers,
	}

	username, password, ok := clientReq.BasicAuth()
	if !ok {
		return nil, unauthResp, nil
	}

	for _, id := range a.params.Identities {
		if id.Username != username {
			continue
		}

		switch id.Hash {
		case hashPlain:
			if id.Password != password {
				continue
			}
		case hashBcrypt:
			if bcrypt.CompareHashAndPassword([]byte(id.Password), []byte(password)) != nil {
				continue
			}
		case hashSHA256:
			if id.Password != fmt.Sprintf("%x", sha256.Sum256([]byte(password))) {
				continue
			}
		default:
			continue
		}

		if a.params.ForwardIdentity {
			clientReq.Header.Set("X-Courier-Authenticated-User", id.Username)
		}

		metrics.Client.Incr(
			metricAuthenticationEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "allow"}),
		)

		return clientReq, nil, nil
	}

	metrics.Client.Incr(
		metricAuthenticationEvaluate,
		util.MergeMaps(tags, map[string]interface{}{"action": "deny"}),
	)

	zap.L().Debug(
		"basic authentication exhausted identity list without valid credentials",
		zap.String("realm", a.params.Realm),
		zap.String("username", username),
	)

	return nil, unauthResp, nil
}
