package supercharged

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestHTTPInvalid(t *testing.T) {
	t.Parallel()

	cases := []*Request{
		// No base URL
		{},
		// Invalid HTTP method
		{
			BaseURL: &url.URL{Scheme: "https", Host: "example.com"},
			Method:  "invalid",
		},
		// Invalid URL
		{
			BaseURL:  &url.URL{Scheme: "https", Host: "example.com"},
			Method:   http.MethodGet,
			Endpoint: "\n", // Control characters disallowed in URLs
		},
		// Non-serializable data payload
		{
			BaseURL: &url.URL{Scheme: "https", Host: "example.com"},
			Method:  http.MethodGet,
			Data:    make(chan struct{}),
		},
	}

	for _, scReq := range cases {
		httpReq, err := scReq.HTTPRequest()

		assert.Error(t, err)
		assert.Nil(t, httpReq)
	}
}

func TestRequestHTTPHeaderData(t *testing.T) {
	t.Parallel()

	scReq := &Request{
		BaseURL:  &url.URL{Scheme: "https", Host: "example.com"},
		Client:   "client",
		Version:  "version",
		Method:   http.MethodGet,
		Endpoint: "/endpoint",
		Data: struct {
			Foo int `json:"foo"`
		}{4},
	}

	httpReq, err := scReq.HTTPRequest()

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/endpoint", httpReq.URL.String())
	assert.Equal(t, http.MethodGet, httpReq.Method)
	assert.Equal(t, "client", httpReq.Header.Get(requestHeaderClientID))
	assert.Equal(t, "version", httpReq.Header.Get(requestHeaderClientVersion))
	assert.Equal(t, "proton/1.0 (client:client; version:version)", httpReq.Header.Get("User-Agent"))
	assert.Equal(t, "{\"foo\":4}", httpReq.Header.Get(requestHeaderData))
	assert.Equal(t, int64(0), httpReq.ContentLength)
}

func TestRequestHTTPBodyData(t *testing.T) {
	t.Parallel()

	scReq := &Request{
		BaseURL:  &url.URL{Scheme: "https", Host: "example.com"},
		Client:   "client",
		Version:  "version",
		Method:   http.MethodPost,
		Endpoint: "/endpoint",
		Data: struct {
			Foo int `json:"foo"`
		}{4},
	}

	httpReq, err := scReq.HTTPRequest()

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/endpoint", httpReq.URL.String())
	assert.Equal(t, http.MethodPost, httpReq.Method)
	assert.Equal(t, "client", httpReq.Header.Get(requestHeaderClientID))
	assert.Equal(t, "version", httpReq.Header.Get(requestHeaderClientVersion))
	assert.Equal(t, "proton/1.0 (client:client; version:version)", httpReq.Header.Get("User-Agent"))
	httpBody, _ := io.ReadAll(httpReq.Body)
	assert.Equal(t, "{\"foo\":4}", string(httpBody))
	assert.Equal(t, int64(9), httpReq.ContentLength)
}

func TestRequestWebsocketURLInvalid(t *testing.T) {
	t.Parallel()

	cases := []*Request{
		// No base URL
		{},
		// Invalid URL scheme
		{
			BaseURL: &url.URL{Scheme: "ftp", Host: "example.com"},
		},
		// Non-serializable data payload
		{
			BaseURL: &url.URL{Scheme: "https", Host: "example.com"},
			Data:    make(chan struct{}),
		},
	}

	for _, scReq := range cases {
		wsURL, err := scReq.WebsocketURL()

		assert.Error(t, err)
		assert.Nil(t, wsURL)
	}
}

func TestRequestWebsocketURLNoData(t *testing.T) {
	t.Parallel()

	scReq := &Request{
		BaseURL:  &url.URL{Scheme: "https", Host: "example.com"},
		Client:   "client",
		Version:  "version",
		Endpoint: "/endpoint",
	}

	wsURL, err := scReq.WebsocketURL()

	assert.NoError(t, err)
	assert.Equal(t, "wss://example.com/endpoint?client=client&version=version", wsURL.String())
}

func TestRequestWebsocketURLData(t *testing.T) {
	t.Parallel()

	scReq := &Request{
		BaseURL:  &url.URL{Scheme: "https", Host: "example.com"},
		Client:   "client",
		Version:  "version",
		Endpoint: "/endpoint",
		Data: struct {
			Foo int `json:"foo"`
		}{4},
	}

	wsURL, err := scReq.WebsocketURL()

	assert.NoError(t, err)
	// Query string serialization order is not guaranteed to be deterministic, so assert on the
	// individual components rather than the serialized URL
	assert.Equal(t, "wss", wsURL.Scheme)
	assert.Equal(t, "example.com", wsURL.Host)
	assert.Equal(t, "/endpoint", wsURL.Path)
	assert.Equal(t, "client", wsURL.Query().Get("client"))
	assert.Equal(t, "version", wsURL.Query().Get("version"))
	assert.Equal(t, "eyJmb28iOjR9", wsURL.Query().Get("data"))
}

func TestParseResponseInvalid(t *testing.T) {
	t.Parallel()

	body := "invalid"
	resp, err := ParseResponse([]byte(body))

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestParseResponseValid(t *testing.T) {
	t.Parallel()

	type payload struct {
		Foo int `json:"foo"`
	}

	var data payload

	body := "{\"success\": true, \"code\": 2, \"message\": \"message\", \"data\": {\"foo\": 4}}"
	resp, err := ParseResponse([]byte(body))

	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
	assert.Equal(t, 2, resp.Code)

	err = json.Unmarshal(resp.Data, &data)

	assert.NoError(t, err)
	assert.Equal(t, 4, data.Foo)
}
