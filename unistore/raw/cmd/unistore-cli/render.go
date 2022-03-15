package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"unistore/schemas/common"
	"unistore/schemas/service"
)

// List of known console output formats.
const (
	outputFormatJSON    = "json"
	outputFormatText    = "text"
	outputFormatYAML    = "yaml"
	outputFormatHuman   = "human"
	outputFormatDefault = ""
)

const (
	// Default rendering indentation level; four spaces.
	defaultIndent = "    "
)

var (
	// errUnsupportedMessage is returned when attempting to use a format renderer with no
	// implementation for the accepted message type.
	errUnsupportedMessage = errors.New("unsupported message type for format renderer")
)

// renderer provides the capability to render formatted messages to an output.
type renderer interface {
	// proto is used for rendering Protobuf messages.
	proto(message proto.Message) error

	// any is used for interface{} messages of arbitrary type.
	any(message interface{}) error
}

// newRenderer creates a new renderer for the desired output format.
func newRenderer(format string, output io.Writer) renderer {
	switch format {
	case outputFormatJSON:
		return &jsonRenderer{output}
	case outputFormatYAML:
		return &yamlRenderer{output}
	case outputFormatHuman:
		return &humanRenderer{output}
	case outputFormatDefault, outputFormatText:
		return &textRenderer{output}
	default:
		return &defaultRenderer{output}
	}
}

// defaultRenderer is a default implementation for unknown output formats.
type defaultRenderer struct {
	output io.Writer
}

// proto errors with errUnsupportedMessage.
func (d *defaultRenderer) proto(message proto.Message) error {
	return errUnsupportedMessage
}

// any errors with errUnsupportedMessage.
func (d *defaultRenderer) any(message interface{}) error {
	return errUnsupportedMessage
}

// jsonRenderer is a JSON renderer.
type jsonRenderer struct {
	output io.Writer
}

// proto reuses the implementation for interface{} types.
func (j *jsonRenderer) proto(message proto.Message) error {
	return j.any(message)
}

// any marshals any type implementing json.Marshaler or any struct with JSON tags.
func (j *jsonRenderer) any(message interface{}) error {
	rendered, err := json.MarshalIndent(message, "", defaultIndent)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(j.output, strings.TrimSpace(string(rendered)))
	return err
}

// yamlRenderer is a YAML renderer.
type yamlRenderer struct {
	output io.Writer
}

// proto errors with errUnsupportedMessage.
// Go-generated code from Protobuf message do not include YAML tags.
func (y *yamlRenderer) proto(message proto.Message) error {
	return errUnsupportedMessage
}

// any marshals any type implementing yaml.Marshaler or any struct with YAML tags.
func (y *yamlRenderer) any(message interface{}) error {
	encoder := yaml.NewEncoder(y.output)
	encoder.SetIndent(len(defaultIndent))
	return encoder.Encode(message)
}

// textRenderer is a Protobuf textproto format renderer.
type textRenderer struct {
	output io.Writer
}

// proto renders the message with the textproto format.
func (t *textRenderer) proto(message proto.Message) error {
	opts := prototext.MarshalOptions{
		Multiline: true,
		Indent:    defaultIndent,
	}

	rendered, err := opts.Marshal(message)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(t.output, strings.TrimSpace(string(rendered)))
	return err
}

// any errors with errUnsupportedMessage.
func (t *textRenderer) any(message interface{}) error {
	return errUnsupportedMessage
}

// humanRenderer is a message type-dependent renderer expressing messages in an easily
// human-readable form. Implementations are defined manually by message type.
type humanRenderer struct {
	output io.Writer
}

// proto renders a human-readable representation of supported Protobuf messages.
func (h *humanRenderer) proto(message proto.Message) error {
	switch m := message.(type) {
	case *service.InfoResponse:
		return h.protoInfoResponse(m)
	case *service.ListBucketsResponse:
		return h.protoListBucketsResponse(m)
	case *service.ListObjectsResponse:
		return h.protoListObjectsResponse(m)
	default:
		return errUnsupportedMessage
	}
}

// protoInfoResponse renders a service.InfoResponse.
func (h *humanRenderer) protoInfoResponse(response *service.InfoResponse) error {
	var renderNode func(node *common.Node) []string

	renderNode = func(node *common.Node) []string {
		var out []string

		childrenConstants := make(map[string]string, len(node.Children))
		childrenNodes := make([]string, 0, len(node.Children))

		for param, value := range node.Children {
			switch c := value.Child.(type) {
			case *common.Node_Value_Value:
				childrenConstants[param] = c.Value
			case *common.Node_Value_Node:
				childrenNodes = append(childrenNodes, param)
			}
		}

		sort.Strings(childrenNodes)

		// Render the node name, along with any constant parameters.
		if len(childrenConstants) > 0 {
			var params []string
			for param, value := range childrenConstants {
				params = append(params, fmt.Sprintf("%s: %s", param, value))
			}
			sort.Strings(params)

			out = append(out, fmt.Sprintf("%s [%s]", node.Name, strings.Join(params, ", ")))
		} else {
			out = append(out, node.Name)
		}

		// Recursively render all children nodes, sorted by name.
		for _, name := range childrenNodes {
			for _, line := range renderNode(node.Children[name].GetNode()) {
				out = append(out, fmt.Sprintf("%s%s", defaultIndent, line))
			}
		}

		return out
	}

	fmt.Fprintln(h.output, "server version:")
	fmt.Fprintf(h.output, "%sunistore/%s\n\n", defaultIndent, response.Version)
	fmt.Fprintln(h.output, "server listener:")
	fmt.Fprintf(h.output, "%s%s\n\n", defaultIndent, response.Address)
	fmt.Fprintln(h.output, "server backend composition tree:")
	for _, line := range renderNode(response.Backend) {
		fmt.Fprintf(h.output, "%s%s\n", defaultIndent, line)
	}

	return nil
}

// protoListBucketsResponse renders a service.ListBucketsResponse.
func (h *humanRenderer) protoListBucketsResponse(response *service.ListBucketsResponse) error {
	var objects service.ListObjectsResponse

	// Reuse the existing renderer for ListObjects by reshaping the ListBuckets response to
	// a list of objects with Object.DIRECTORY type.

	for _, bucket := range response.Buckets {
		head := &service.HeadObjectResponse{
			Metadata: &common.Metadata{
				Key:          bucket.Bucket.Name,
				ModifiedTime: bucket.Bucket.Timestamp,
				ObjectType:   common.Object_DIRECTORY,
			},
		}

		objects.Objects = append(objects.Objects, head)
	}

	return h.protoListObjectsResponse(&objects)
}

// protoListObjectsResponse renders a service.ListObjectsResponse.
func (h *humanRenderer) protoListObjectsResponse(response *service.ListObjectsResponse) error {
	// Sort the objects alphabetically by key name, while grouping directories at the
	// beginning of the list.
	sort.Slice(response.Objects, func(i int, j int) bool {
		l := response.Objects[i]
		r := response.Objects[j]

		if l.Metadata.ObjectType == common.Object_DIRECTORY &&
			r.Metadata.ObjectType == common.Object_DIRECTORY {
			// Both objects are directories; sort alphabetically
			return l.Metadata.Key < r.Metadata.Key
		} else if l.Metadata.ObjectType == common.Object_DIRECTORY {
			// The left object is a directory; always prefer it
			return true
		} else if r.Metadata.ObjectType == common.Object_DIRECTORY {
			// The right object is a directory; always prefer it
			return false
		} else {
			// Both objects are non-directories; sort alphabetically
			return l.Metadata.Key < r.Metadata.Key
		}
	})

	for _, object := range response.Objects {
		switch object.Metadata.ObjectType {
		case common.Object_DIRECTORY:
			if !strings.HasSuffix(object.Metadata.Key, "/") {
				object.Metadata.Key = fmt.Sprintf("%s/", object.Metadata.Key)
			}

			_, err := fmt.Fprintf(
				h.output,
				"%v\t%s\t%s\n",
				object.Metadata.ModifiedTime.AsTime().Format(time.RFC3339),
				"-",
				object.Metadata.Key,
			)
			if err != nil {
				return err
			}

		case common.Object_REGULAR:
			_, err := fmt.Fprintf(
				h.output,
				"%v\t%s\t%s\n",
				object.Metadata.ModifiedTime.AsTime().Format(time.RFC3339),
				humanize.Bytes(object.Metadata.Size),
				object.Metadata.Key,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// any errors with errUnsupportedMessage.
func (h *humanRenderer) any(message interface{}) error {
	return errUnsupportedMessage
}
