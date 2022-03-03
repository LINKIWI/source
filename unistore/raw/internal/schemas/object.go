package schemas

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"unistore/schemas/common"
)

// Canonical attribute names for known attribute types.
const (
	attrPrefixMetadata = "unistore.metadata"
	attrPrefixExtra    = "unistore.extra"

	attrServer      = attrPrefixMetadata + ".server"
	attrVersion     = attrPrefixMetadata + ".version"
	attrTimestamp   = attrPrefixMetadata + ".timestamp"
	attrBackend     = attrPrefixMetadata + ".backend"
	attrCompression = attrPrefixMetadata + ".compression"
	attrEncryption  = attrPrefixMetadata + ".encryption"
	attrSize        = attrPrefixMetadata + ".size"
)

var (
	// extraAttrRe is a regular expression denoting acceptable characters in an attribute name.
	extraAttrRe = regexp.MustCompile("^[a-z0-9-_]{1,128}$")
)

// AttributesEncodingOptions describes options for marshaling and unmarshalling attributes.
type AttributesEncodingOptions struct {
	Prefix string
}

// MarshalAttributes marshals a structured set of attributes to a key-value string map suitable for
// storage in the underlying backend.
func MarshalAttributes(attributes *common.Attributes, opts *AttributesEncodingOptions) (map[string]string, error) {
	if opts == nil {
		opts = &AttributesEncodingOptions{}
	}

	out := map[string]string{
		attrServer:      attributes.Server,
		attrVersion:     attributes.Version,
		attrTimestamp:   attributes.Timestamp.AsTime().Format(time.RFC3339Nano),
		attrBackend:     attributes.Backend.String(),
		attrCompression: attributes.Compression.String(),
		attrEncryption:  attributes.Encryption.String(),
		attrSize:        fmt.Sprintf("%d", attributes.Size),
	}

	if attributes.Extra != nil {
		for extra, value := range attributes.Extra {
			if !extraAttrRe.MatchString(extra) {
				return nil, fmt.Errorf(
					"schemas: attribute name is illegal: pattern=%v attribute=%s",
					extraAttrRe,
					extra,
				)
			}

			out[fmt.Sprintf("%s.%s", attrPrefixExtra, extra)] = value
		}
	}

	if opts.Prefix != "" {
		prefixed := make(map[string]string)

		for attr, value := range out {
			prefixed[fmt.Sprintf("%s.%s", opts.Prefix, attr)] = value
		}

		out = prefixed
	}

	return out, nil
}

// UnmarshalAttributes parses a key-value string map of attributes into a structured data.
func UnmarshalAttributes(attributes map[string]string, opts *AttributesEncodingOptions) (*common.Attributes, error) {
	if opts == nil {
		opts = &AttributesEncodingOptions{}
	}

	out := &common.Attributes{
		Extra: make(map[string]string),
	}

	// Each element is keyed by the canonical attribute name, and points to a function that
	// can arbitrarily modify a *common.Attributes pointer.
	parsers := map[string]func(attr string, value string, out *common.Attributes) error{
		attrServer: func(attr string, value string, out *common.Attributes) error {
			out.Server = value
			return nil
		},
		attrVersion: func(attr string, value string, out *common.Attributes) error {
			out.Version = value
			return nil
		},
		attrTimestamp: func(attr string, value string, out *common.Attributes) error {
			t, err := time.Parse(time.RFC3339Nano, value)
			if err != nil {
				return fmt.Errorf(
					"schemas: error parsing timestamp: timestamp=%s err=%v",
					value,
					err,
				)
			}

			out.Timestamp = timestamppb.New(t)
			return nil
		},
		attrBackend: func(attr string, value string, out *common.Attributes) error {
			if value == "" {
				out.Backend = common.Backend_EMPTY
				return nil
			}

			if c, ok := common.Backend_value[strings.ToUpper(value)]; ok {
				out.Backend = common.Backend(c)
				return nil
			}

			return fmt.Errorf("schemas: unknown backend value: backend=%s", value)
		},
		attrCompression: func(attr string, value string, out *common.Attributes) error {
			if c, ok := common.Compression_value[strings.ToUpper(value)]; ok {
				out.Compression = common.Compression(c)
				return nil
			}

			return fmt.Errorf("schemas: unknown compression value: compression=%s", value)
		},
		attrEncryption: func(attr string, value string, out *common.Attributes) error {
			if e, ok := common.Encryption_value[strings.ToUpper(value)]; ok {
				out.Encryption = common.Encryption(e)
				return nil
			}

			return fmt.Errorf("schemas: unknown encryption value: encryption=%s", value)
		},
		attrSize: func(attr string, value string, out *common.Attributes) error {
			s, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf(
					"schemas: error parsing object size: size=%s err=%v",
					value,
					err,
				)
			}

			out.Size = s
			return nil
		},
		// Default, catch-all visitor
		"": func(attr string, value string, out *common.Attributes) error {
			// Custom extra attributes
			if strings.HasPrefix(attr, fmt.Sprintf("%s.", attrPrefixExtra)) {
				out.Extra[attr[len(attrPrefixExtra)+1:]] = value
			}

			return nil
		},
	}

	for attr, value := range attributes {
		if opts.Prefix != "" && strings.HasPrefix(attr, fmt.Sprintf("%s.", opts.Prefix)) {
			attr = attr[len(opts.Prefix)+1:]
		}

		parser, ok := parsers[attr]
		if !ok {
			parser = parsers[""]
		}

		if err := parser(attr, value, out); err != nil {
			return nil, err
		}
	}

	return out, nil
}
