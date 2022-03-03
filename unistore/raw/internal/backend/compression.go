package backend

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zstd"

	"unistore/internal/schemas"
	"unistore/internal/util"
	"unistore/schemas/common"
)

// Compression is a Backend that performs in-flight compression and decompression of get and put
// data payloads.
type Compression struct {
	algorithm common.Compression
	level     int
	Backend
}

// NewCompression creates a new Compression with the specified algorithm and compression level.
func NewCompression(algorithm string, level int, backend Backend) Backend {
	var compressor func(io.Writer) (io.WriteCloser, error)
	var decompressor func(io.Reader) (io.ReadCloser, error)

	parsed := common.Compression(common.Compression_value[strings.ToUpper(algorithm)])

	switch parsed {
	case common.Compression_GZIP:
		compressor = func(raw io.Writer) (io.WriteCloser, error) {
			if level == 0 {
				level = gzip.DefaultCompression
			}

			return gzip.NewWriterLevel(raw, level)
		}
		decompressor = func(raw io.Reader) (io.ReadCloser, error) {
			return gzip.NewReader(raw)
		}

	case common.Compression_ZSTD:
		compressor = func(raw io.Writer) (io.WriteCloser, error) {
			return zstd.NewWriter(
				raw,
				zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(level)),
			)
		}
		decompressor = func(raw io.Reader) (io.ReadCloser, error) {
			decoder, err := zstd.NewReader(raw)
			if err != nil {
				return nil, err
			}

			return io.NopCloser(decoder), nil
		}

	case common.Compression_SNAPPY:
		compressor = func(raw io.Writer) (io.WriteCloser, error) {
			var opts []s2.WriterOption

			switch level {
			case 0:
				opts = append(opts, s2.WriterUncompressed())
			case 1:
				// Default level is fast
			case 2:
				opts = append(opts, s2.WriterBetterCompression())
			case 3:
				opts = append(opts, s2.WriterBestCompression())
			default:
				return nil, fmt.Errorf(
					"backend: unknown snappy compression level: level=%d",
					level,
				)
			}

			return s2.NewWriter(raw, opts...), nil
		}
		decompressor = func(raw io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(s2.NewReader(raw)), nil
		}

	case common.Compression_PLAINTEXT:
		compressor = func(raw io.Writer) (io.WriteCloser, error) {
			return util.NopWriteCloser(raw), nil
		}
		decompressor = func(raw io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(raw), nil
		}
	}

	c := &Compression{
		algorithm: parsed,
		level:     level,
		Backend:   newIOProcessor(decompressor, compressor, 0, 0, backend),
	}

	return newInstrumentation("compression", c)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (c *Compression) Descriptor() *common.Node {
	return &common.Node{
		Name: "compression",
		Children: map[string]*common.Node_Value{
			"algorithm": {
				Child: &common.Node_Value_Value{
					Value: c.algorithm.String(),
				},
			},
			"level": {
				Child: &common.Node_Value_Value{
					Value: fmt.Sprintf("%d", c.level),
				},
			},
			"backend": {
				Child: &common.Node_Value_Node{
					Node: c.Backend.Descriptor(),
				},
			},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (c *Compression) String() string {
	return schemas.MarshalDescriptor(c.Descriptor())
}
