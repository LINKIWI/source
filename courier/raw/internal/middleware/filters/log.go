package filters

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"lib.kevinlin.info/aperture/lib"

	"courier/internal/config"
	"courier/internal/meta"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricLogLineWrite = "filter.log.line.write"
	metricLogLineSize  = "filter.log.line.size"
)

const (
	logSerializerJSON = "json"
	logSerializerLine = "line"
)

// jsonLogEntry is an alias type for a log entry that can be serialized as JSON.
type jsonLogEntry logEntry

// lineLogEntry is an alias type for a log entry that can be serialized as a formatted line.
type lineLogEntry logEntry

// logEntry is a structured container for information relevant in a log trace.
type logEntry struct {
	Timestamp     time.Time         `json:"timestamp" format:"timestamp"`
	RequestID     uuid.UUID         `json:"request_id" format:"request_id"`
	Listener      string            `json:"listener" format:"listener"`
	Transport     string            `json:"transport" format:"transport"`
	Protocol      string            `json:"protocol" format:"protocol"`
	Scheme        string            `json:"scheme" format:"scheme"`
	URL           string            `json:"url" format:"url"`
	ProxyURL      string            `json:"proxy_url" format:"proxy_url"`
	Method        string            `json:"method" format:"method"`
	Host          string            `json:"host" format:"host"`
	Path          string            `json:"path" format:"path"`
	StatusCode    int               `json:"status" format:"status"`
	RemoteIP      string            `json:"remote_ip" format:"remote_ip"`
	RequestSize   int64             `json:"request_size" format:"request_size"`
	ResponseSize  int64             `json:"response_size" format:"response_size"`
	UserAgent     string            `json:"user_agent" format:"user_agent"`
	Duration      float64           `json:"duration" format:"duration"`
	TLSServerName string            `json:"tls_server_name" format:"tls_server_name"`
	TLSClientName string            `json:"tls_client_name" format:"tls_client_name"`
	TLSVersion    string            `json:"tls_version" format:"tls_version"`
	Version       string            `json:"version" format:"version"`
	Tags          map[string]string `json:"tags"`
}

// logParams is the configuration descriptor for the log filter.
type logParams struct {
	Match      *matchSpec `yaml:"match"`
	Path       string     `yaml:"path"`
	Serializer string     `yaml:"serializer"`
	JSON       struct{}   `yaml:"json"`
	Line       struct {
		Format string `yaml:"format"`
	} `yaml:"line"`
	Tags []struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"tags"`
}

// log is a filter that logs completed requests to a file on disk.
type log struct {
	file   *os.File
	params *logParams
}

// NewLog creates a new log filter.
func NewLog(cfg *config.Filter) (middleware.Filter, error) {
	var params logParams
	var filter middleware.Filter

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse log filter params",
			StackedError: err,
		}
	} else if params.Path == "" {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "missing log file path",
		}
	} else if len(params.Serializer) > 0 &&
		params.Serializer != logSerializerJSON &&
		params.Serializer != logSerializerLine {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "unknown log serializer",
			Tags: map[string]interface{}{
				"serializer": params.Serializer,
			},
		}
	}

	file, err := os.OpenFile(params.Path, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, 0664)
	if err != nil {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "failed to open log file for writing",
			Tags: map[string]interface{}{
				"path": params.Path,
			},
			StackedError: err,
		}
	}

	filter = &log{file: file, params: &params}
	if params.Match != nil {
		filter = newMatch("log", params.Match, filter)
	}

	return middleware.NewAsyncFilter(middleware.NewInstrumentedFilter("log", filter)), nil
}

// ProcessRequest assigns the request a unique ID and starts a latency stopwatch.
func (l *log) ProcessRequest(clientReq *http.Request) (*http.Request, *http.Response, error) {
	ctx := context.WithValue(clientReq.Context(), ctxLogRequestID, uuid.New())
	ctx = context.WithValue(ctx, ctxLogStopwatch, lib.NewStopwatch())
	ctx = context.WithValue(ctx, ctxLogOriginalURL, clientReq.URL)

	return clientReq.WithContext(ctx), nil, nil
}

// ProcessResponse assembles information about the request/response pair and logs the entry.
func (l *log) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	var serialized []byte
	var transport string

	logTags := make(map[string]string)
	metricTags := map[string]interface{}{
		"route_host": proxyResp.Request.Host,
		"log_path":   l.params.Path,
		"serializer": l.params.Serializer,
	}

	listener := proxyResp.Request.Context().Value(http.LocalAddrContextKey).(net.Addr)
	requestID := proxyResp.Request.Context().Value(ctxLogRequestID).(uuid.UUID)
	stopwatch := proxyResp.Request.Context().Value(ctxLogStopwatch).(*lib.Stopwatch)
	originalURL := proxyResp.Request.Context().Value(ctxLogOriginalURL).(*url.URL)

	remoteIP, _, err := net.SplitHostPort(proxyResp.Request.RemoteAddr)
	if err != nil {
		if proxyResp.Request.RemoteAddr == "@" {
			remoteIP = "@" // Special address for Unix socket transports
		} else {
			return proxyResp, &util.Error{
				Namespace:    "filters",
				Message:      "unable to parse client remote address",
				Tags:         map[string]interface{}{"address": proxyResp.Request.RemoteAddr},
				StackedError: err,
			}
		}
	}

	switch addr := listener.(type) {
	case *net.TCPAddr:
		if addr.IP.To4() != nil {
			transport = "tcp4"
		} else {
			transport = "tcp6"
		}
	case *net.UnixAddr:
		transport = "unix"
	default:
		transport = "unknown"
	}

	for _, tag := range l.params.Tags {
		logTags[tag.Key] = tag.Value
	}

	entry := logEntry{
		Timestamp:    time.Now(),
		RequestID:    requestID,
		Listener:     listener.String(),
		Transport:    transport,
		URL:          proxyResp.Request.RequestURI,
		ProxyURL:     proxyResp.Request.URL.String(),
		Protocol:     proxyResp.Request.Proto,
		Scheme:       originalURL.Scheme,
		Method:       proxyResp.Request.Method,
		Host:         originalURL.Host,
		Path:         originalURL.Path,
		StatusCode:   proxyResp.StatusCode,
		RemoteIP:     remoteIP,
		RequestSize:  proxyResp.Request.ContentLength,
		ResponseSize: proxyResp.ContentLength,
		UserAgent:    proxyResp.Request.UserAgent(),
		Duration:     float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Version:      meta.Version,
		Tags:         logTags,
	}

	if tlsConnState := proxyResp.Request.TLS; tlsConnState != nil {
		entry.TLSServerName = tlsConnState.ServerName

		if len(tlsConnState.PeerCertificates) > 0 {
			entry.TLSClientName = strings.Join(
				tlsConnState.PeerCertificates[0].DNSNames,
				",",
			)
		}

		switch tlsConnState.Version {
		case tls.VersionTLS11:
			entry.TLSVersion = "TLSv1.1"
		case tls.VersionTLS12:
			entry.TLSVersion = "TLSv1.2"
		case tls.VersionTLS13:
			entry.TLSVersion = "TLSv1.3"
		default:
			entry.TLSVersion = "unknown"
		}
	}

	switch l.params.Serializer {
	case logSerializerLine:
		serialized, err = lineLogEntry(entry).serialize(l.params.Line.Format)
	default:
		serialized, err = jsonLogEntry(entry).serialize()
	}

	if err != nil {
		return proxyResp, &util.Error{
			Namespace:    "filters",
			Message:      "failed to serialize log line",
			StackedError: err,
		}
	}

	if _, err := l.file.Write(append(serialized, []byte("\n")...)); err != nil {
		return proxyResp, &util.Error{
			Namespace:    "filters",
			Message:      "failed to write serialized log line to file",
			StackedError: err,
		}
	}

	metrics.Client.Incr(metricLogLineWrite, metricTags)
	metrics.Client.Size(metricLogLineSize, int64(len(serialized)), metricTags)

	return proxyResp, nil
}

// serialize formats the log trace in accordance with the config-specified format string, replacing
// struct values with corresponding tag keys for tag `format`.
func (ll lineLogEntry) serialize(format string) ([]byte, error) {
	var replace []string

	t := reflect.TypeOf(ll)
	v := reflect.ValueOf(ll)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		formatKey, ok := field.Tag.Lookup("format")
		if !ok {
			continue
		}

		replace = append(
			replace,
			fmt.Sprintf("{%s}", formatKey),
			fmt.Sprintf("%v", v.Field(i).Interface()),
		)
	}

	// Add custom tags to the formatter as well, whose format strings are prefixed with tag_
	for key, value := range ll.Tags {
		replace = append(replace, fmt.Sprintf("{tag_%s}", key), value)
	}

	return []byte(strings.NewReplacer(replace...).Replace(format)), nil
}

// serialize serializes the log entry as JSON.
func (jl jsonLogEntry) serialize() ([]byte, error) {
	return json.Marshal(&jl)
}
