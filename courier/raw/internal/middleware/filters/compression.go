package filters

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"sort"

	headerutil "github.com/golang/gddo/httputil/header"
	"go.uber.org/zap"

	"courier/internal/config"
	"courier/internal/metrics"
	"courier/internal/middleware"
	"courier/internal/util"
)

const (
	metricCompressionAlgorithmNegotiate   = "filter.compression.negotiate"
	metricCompressionAlgorithmUnavailable = "filter.compression.unavailable"
	metricCompressionAlgorithmSkip        = "filter.compression.skip"
)

const (
	algorithmWildcard = "*"
	algorithmIdentity = "identity"
	algorithmGzip     = "gzip"
	algorithmDeflate  = "deflate"
)

// compressorFactory is a factory for generating a compressor that accepts raw writes and performs
// streaming compression while writing to the underlying io.Writer.
type compressorFactory func(raw io.Writer) (io.WriteCloser, error)

// compressionParams is the configuration descriptor for the compression filter.
type compressionParams struct {
	Match      *matchSpec `yaml:"match"`
	Algorithms []struct {
		Name  string `yaml:"name"`
		Level int    `yaml:"level"`
	} `yaml:"algorithms"`
}

// compression applies opt-in compression to the outgoing response.
type compression struct {
	params      *compressionParams
	compressors map[string]compressorFactory
	noop
}

// NewCompression creates a new compression filter.
func NewCompression(cfg *config.Filter) (middleware.Filter, error) {
	var params compressionParams
	var filter middleware.Filter
	compressors := make(map[string]compressorFactory)

	if err := cfg.Params.Decode(&params); err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "failed to parse compression filter params",
			StackedError: err,
		}
	} else if len(params.Algorithms) == 0 {
		return nil, &util.Error{
			Namespace: "filters",
			Message:   "at least one compression algorithm must be enabled",
		}
	}

	for _, algorithm := range params.Algorithms {
		if algorithm.Level == 0 {
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "compression level must be nonzero (use -1 for default)",
			}
		}

		switch algorithm.Name {
		case "gzip":
			compressors[algorithmGzip] = func(writer io.Writer) (io.WriteCloser, error) {
				return gzip.NewWriterLevel(writer, algorithm.Level)
			}
		case "deflate", "zlib":
			compressors[algorithmDeflate] = func(writer io.Writer) (io.WriteCloser, error) {
				return zlib.NewWriterLevel(writer, algorithm.Level)
			}
		case "":
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "compression algorithm name missing",
			}
		default:
			return nil, &util.Error{
				Namespace: "filters",
				Message:   "unknown compression algorithm",
				Tags:      map[string]interface{}{"algorithm": algorithm},
			}
		}
	}

	filter = &compression{compressors: compressors, params: &params}
	if params.Match != nil {
		filter = newMatch("compression", params.Match, filter)
	}

	return middleware.NewInstrumentedFilter("compression", filter), nil
}

// ProcessResponse applies streaming compression of the proxy response contingent on the compression
// algorithms supported by the client as advertised in the incoming Accept-Encoding header.
func (c *compression) ProcessResponse(proxyResp *http.Response) (*http.Response, error) {
	tags := map[string]interface{}{
		"route_host": proxyResp.Request.Host,
		"method":     proxyResp.Request.Method,
		"protocol":   proxyResp.Request.Proto,
	}

	// Response compression is not supported for bidirectional streams, e.g. websockets or
	// HTTP/2 streams.
	if _, ok := proxyResp.Body.(io.Writer); ok {
		metrics.Client.Incr(metricCompressionAlgorithmSkip, tags)
		return proxyResp, nil
	}

	// Don't modify the response if the upstream has itself compressed the response, evidenced
	// by a non-empty content encoding.
	if proxyResp.Header.Get("Content-Encoding") != "" {
		metrics.Client.Incr(metricCompressionAlgorithmSkip, tags)
		return proxyResp, nil
	}

	// Parse the spec of encodings accepted by the client and sort in descending order of
	// preference based on the optional quality value (weight).
	encodings := headerutil.ParseAccept(proxyResp.Request.Header, "Accept-Encoding")
	sort.Slice(encodings, func(i int, j int) bool { return encodings[i].Q > encodings[j].Q })

	for _, encoding := range encodings {
		if encoding.Value == algorithmIdentity {
			// Explicitly apply no compression.
			return proxyResp, nil
		} else if encoding.Value == algorithmWildcard {
			// Choose any available compression algorithm.
			for algorithm := range c.compressors {
				encoding.Value = algorithm
				break
			}
		}

		compressor, ok := c.compressors[encoding.Value]
		if !ok {
			// No implementation available for the requested compression.
			continue
		}

		metrics.Client.Incr(
			metricCompressionAlgorithmNegotiate,
			util.MergeMaps(tags, map[string]interface{}{"algorithm": encoding.Value}),
		)
		return c.processCompressedResponse(
			proxyResp,
			encoding.Value,
			compressor,
		)
	}

	if len(encodings) > 0 {
		metrics.Client.Incr(metricCompressionAlgorithmUnavailable, tags)
		zap.L().Warn(
			"no available compression implementation satisfying requested encoding(s)",
			zap.Any("encodings", encodings),
		)
	}

	return proxyResp, nil
}

// processCompressedResponse applies compression to the response body based on the pluggable
// compressor factory.
func (c *compression) processCompressedResponse(
	proxyResp *http.Response,
	name string,
	compressor compressorFactory,
) (*http.Response, error) {
	var compressionWriter io.WriteCloser

	proxyResp.Header.Del("Content-Length") // The original body length is no longer accurate.
	proxyResp.Header.Set("Content-Encoding", name)

	respBody := proxyResp.Body
	proxyResp.Body, compressionWriter = io.Pipe()
	rawWriter, err := compressor(compressionWriter)

	if err != nil {
		return nil, &util.Error{
			Namespace:    "filters",
			Message:      "error creating compression writer",
			StackedError: err,
		}
	}

	go func() {
		// Directly pipe the original response body into the compression writer. The output
		// of the compression writer is then piped directly into a reader for the final
		// compressed response body.
		_, err := io.Copy(rawWriter, respBody)
		if err != nil {
			zap.L().Error("error during compression response pipe", zap.Error(err))
		}

		rawWriter.Close()
		compressionWriter.Close()
		respBody.Close()
	}()

	return proxyResp, nil
}
