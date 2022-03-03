package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// response is a structured container for a generic response.
type response struct {
	status  int
	code    string
	message string
	tags    map[string]string
}

// response creates an http.Response from the response metadata.
func (e *response) response() *http.Response {
	var body []string

	body = append(body, e.message)
	for key, value := range e.tags {
		body = append(body, fmt.Sprintf("%v=%v", key, value))
	}

	return &http.Response{
		StatusCode: e.status,
		Body:       io.NopCloser(strings.NewReader(strings.Join(body, "\n"))),
		Header:     make(http.Header),
	}
}

// String provides a string representation of the response.
func (e *response) String() string {
	return fmt.Sprintf("[%d] %s: %s", e.status, e.code, e.message)
}

// Error provides a string representation of the response.
func (e *response) Error() string {
	return e.String()
}
