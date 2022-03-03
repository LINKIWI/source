package b2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Status is a structured error status returned by the B2 API in the response body.
// Reference: https://www.backblaze.com/b2/docs/calling.html#error_handling
type Status struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface by providing a string representation of the status.
func (s *Status) Error() string {
	return fmt.Sprintf("HTTP %d: %s (code %s)", s.Status, s.Message, s.Code)
}

// statusFromError creates a Status from a vanilla error.
func statusFromError(err error) *Status {
	if status, ok := err.(*Status); ok {
		return status
	}

	return &Status{
		Status:  http.StatusInternalServerError,
		Code:    "unknown",
		Message: err.Error(),
	}
}

// URL wraps a *url.URL while implementing json.Unmarshaler in order to allow parsing a JSON string
// field directly into a URL.
type URL struct {
	*url.URL
}

// UnmarshalJSON attempts to parse the data as a JSON string, followed by parsing it as a URL.
func (u *URL) UnmarshalJSON(data []byte) error {
	var err error
	var raw string

	// First parse the raw data as a string (the expected source type for a *url.URL).
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	u.URL, err = url.Parse(raw)

	return err
}

// UnixTimestamp wraps a time.Time while implementing json.Unmarshaler in order to allow parsing a
// JSON number as a timestamp (in milliseconds).
type UnixTimestamp struct {
	time.Time
}

// UnmarshalJSON attempts to parse the data as an integer, followed by parsing it as a Unix
// timestamp in millisecond units.
func (u *UnixTimestamp) UnmarshalJSON(data []byte) error {
	var raw int64

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert milliseconds to nanoseconds
	u.Time = time.Unix(0, raw*1000*1000)

	return nil
}

// marshalHeaders marshals a struct into a http.Header based on the "header" field tags in the
// struct definition. This is useful for serializing a struct declaration into a set of headers
// available for use in an outgoing http.Request.
//
// The supported tag format is as follows:
//
//     <header key>[,[encode]]
//
// <header key> is a required string literal with the desired header key for the entry; the field
// value is used as the corresponding header value. Note that this implementation uses add semantics
// for headers; if the same header key is defined multiple times in the struct, it will appear
// multiple times in the serialized http.Header.
//
// Following the mandatory header key is an optional comma which delimits the header key and any
// optional boolean toggles. To use all default behaviors, the comma can be elided.
//
// encode is an optional keyword that instructs the serializer to URL-encode the associated header
// value in the serialized header.
//
// For example, the following struct definition:
//
//     type Foo struct {
//             Bar string `header:"X-Bar,encode"`
//	       Baz int    `header:"X-Baz"`
//             Qux int
//     }
//     foo := &Foo{Bar: "hello/world", Baz: 4}
//
// Would be serialized as the following headers:
//
//     X-Bar: hello%2Fworld
//     X-Baz: 4
func marshalHeaders(request interface{}) (http.Header, error) {
	h := make(http.Header)

	t := reflect.TypeOf(request)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("b2: request interface{} is not a struct: kind=%v", t.Kind())
	}

	v := reflect.ValueOf(request)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		raw, ok := field.Tag.Lookup("header")
		if !ok {
			continue
		}

		value := fmt.Sprintf("%v", v.Field(i).Interface())
		if value == "" {
			continue
		}

		spec := strings.Split(raw, ",")
		if len(spec) == 0 {
			return nil, fmt.Errorf("b2: header specification provided with no data")
		} else if len(spec) == 1 && spec[0] == "" {
			return nil, fmt.Errorf("b2: header specification missing header key")
		}

		// Marshal options are specified as keywords following the header key, which is the
		// first component in the comma-delimited specification.
		opts := make(map[string]bool)
		if len(spec) > 1 {
			for i := 1; i < len(spec); i++ {
				opts[spec[i]] = true
			}
		}

		if _, ok := opts["encode"]; ok {
			value = url.QueryEscape(value)
		}

		h.Add(spec[0], value)
	}

	return h, nil
}
