package filters

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricMatchEvaluate = "filter.match.evaluate"
)

// matchSpec is a specification of parameters against which the incoming request should match.
// All fields are optional; specification of no criteria indicates that all requests should be
// matched. Similarly, specification of multiple criteria indicates that all criteria need to be
// met for the request to be matched.
type matchSpec struct {
	Scheme    *config.Regex `yaml:"scheme"`
	Method    *config.Regex `yaml:"method"`
	Host      *config.Regex `yaml:"host"`
	Path      *config.Regex `yaml:"path"`
	Transport *config.Regex `yaml:"transport"`
	Source    *config.CIDR  `yaml:"source"`
	Header    *struct {
		Key   string        `yaml:"key"`
		Value *config.Regex `yaml:"value"`
	} `yaml:"header"`
}

// match is a meta filter that passes logic to another filter only if the request fulfills the match
// specification criteria.
type match struct {
	name string
	spec *matchSpec
	middleware.Filter
}

// newMatch wraps another middleware.Filter with client request match criteria.
func newMatch(name string, spec *matchSpec, filter middleware.Filter) middleware.Filter {
	return &match{
		name:   name,
		spec:   spec,
		Filter: filter,
	}
}

// ProcessRequest returns the request as-is without modification if the request fails to meet all
// the specified criteria, and passes the request to the underlying filter if the request meets
// all the specified criteria.
func (m *match) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	tags := map[string]interface{}{
		"filter":     m.name,
		"route_host": clientReq.Host,
		"method":     clientReq.Method,
		"protocol":   clientReq.Proto,
	}

	ctx := clientReq.Context()

	passthrough, err := m.match(clientReq)
	if err != nil {
		return nil, nil, &util.Error{
			Namespace:    "filters",
			Message:      "error evaluating request against match specification",
			Tags:         map[string]interface{}{"match": m.spec},
			StackedError: err,
		}
	}

	// Maintain request lifecycle-attached state on whether the request matched the
	// specification, and thus was eligible for filter processing.
	status, ok := ctx.Value(ctxMatchStatus).(map[middleware.Filter]bool)
	if ok {
		status[m.Filter] = passthrough
	} else {
		ctx = context.WithValue(
			ctx,
			ctxMatchStatus,
			map[middleware.Filter]bool{m.Filter: passthrough},
		)
	}

	if !passthrough {
		metrics.Client.Incr(
			metricMatchEvaluate,
			util.MergeMaps(tags, map[string]interface{}{"action": "reject"}),
		)
		return clientReq.WithContext(ctx), nil, nil
	}

	metrics.Client.Incr(
		metricMatchEvaluate,
		util.MergeMaps(tags, map[string]interface{}{"action": "pass"}),
	)
	zap.L().Debug(
		"request matched specification criteria",
		zap.Stringer("method", m.spec.Method),
		zap.Stringer("host", m.spec.Host),
		zap.Stringer("path", m.spec.Path),
		zap.Any("header", m.spec.Header),
	)

	return m.Filter.ProcessRequest(clientReq.WithContext(ctx))
}

// ProcessResponse conditionally invokes the underlying filter's response processor if and only if
// the same request was eligible for the filter's request processing.
func (m *match) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	status := proxyResp.Request.Context().Value(ctxMatchStatus).(map[middleware.Filter]bool)

	if passthrough, ok := status[m.Filter]; ok && passthrough {
		return m.Filter.ProcessResponse(proxyResp)
	}

	return proxyResp, nil
}

// match evaluates the incoming client request against the match specification.
func (m *match) match(clientReq *http.Request) (bool, error) {
	if m.spec.Scheme != nil && !m.spec.Scheme.MatchString(clientReq.URL.Scheme) {
		return false, nil
	} else if m.spec.Method != nil && !m.spec.Method.MatchString(clientReq.Method) {
		return false, nil
	} else if m.spec.Host != nil && !m.spec.Host.MatchString(clientReq.Host) {
		return false, nil
	} else if m.spec.Path != nil && !m.spec.Path.MatchString(clientReq.URL.Path) {
		return false, nil
	} else if m.spec.Header != nil {
		value := clientReq.Header.Get(m.spec.Header.Key)
		if value == "" {
			return false, nil
		}

		if m.spec.Header.Value != nil && !m.spec.Header.Value.MatchString(value) {
			return false, nil
		}
	} else if m.spec.Transport != nil {
		if clientIP, _, err := net.SplitHostPort(clientReq.RemoteAddr); err == nil {
			ip := net.ParseIP(clientIP)

			// IPv4 addresses should match either "tcp" or "tcp4"
			if ip.To4() != nil &&
				!m.spec.Transport.MatchString("tcp") &&
				!m.spec.Transport.MatchString("tcp4") {
				return false, nil
			}

			// IPv6 addresses should match either "tcp" or "tcp6"
			if ip.To4() == nil &&
				!m.spec.Transport.MatchString("tcp") &&
				!m.spec.Transport.MatchString("tcp6") {
				return false, nil
			}
		}

		if clientReq.RemoteAddr == "@" && !m.spec.Transport.MatchString("unix") {
			return false, nil
		}
	} else if m.spec.Source != nil {
		clientIP, _, err := net.SplitHostPort(clientReq.RemoteAddr)
		if err != nil {
			return false, err
		}

		if !m.spec.Source.Contains(net.ParseIP(clientIP)) {
			return false, nil
		}
	}

	return true, nil
}
