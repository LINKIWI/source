package proton

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/proton/supercharged"
)

type staticResponseTransport struct {
	resp *http.Response
}

func (t *staticResponseTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return t.resp, nil
}

func TestClientDoHTTPProtocolError(t *testing.T) {
	t.Parallel()

	transport := &staticResponseTransport{
		resp: &http.Response{
			Body: io.NopCloser(strings.NewReader("foo")),
		},
	}

	cfg, _ := NewConfig(&Config{
		BaseURL: &url.URL{Scheme: "http", Host: "example.com"},
		Backend: &http.Client{Transport: transport},
	})
	client := NewClient(cfg)

	err := client.DoHTTP(http.MethodGet, "/", nil, nil)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, supercharged.CodeClientUndefined, err.Code)
}

func TestClientDoHTTPSuperchargedError(t *testing.T) {
	t.Parallel()

	resp, _ := json.Marshal(&supercharged.Response{
		Success: false,
		Code:    5,
		Message: "message",
	})

	transport := &staticResponseTransport{
		resp: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader(resp)),
		},
	}

	cfg, _ := NewConfig(&Config{
		BaseURL: &url.URL{Scheme: "http", Host: "example.com"},
		Backend: &http.Client{Transport: transport},
	})
	client := NewClient(cfg)

	err := client.DoHTTP(http.MethodGet, "/", nil, nil)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, 5, err.Code)
	assert.Equal(t, "message", err.Message)
}

func TestClientDoHTTPDataParse(t *testing.T) {
	t.Parallel()

	type payload struct {
		Foo int `json:"foo"`
	}

	data, _ := json.Marshal(&payload{Foo: 5})
	resp, _ := json.Marshal(&supercharged.Response{
		Success: true,
		Data:    data,
	})

	transport := &staticResponseTransport{
		resp: &http.Response{
			Body: io.NopCloser(bytes.NewReader(resp)),
		},
	}

	cfg, _ := NewConfig(&Config{
		BaseURL: &url.URL{Scheme: "http", Host: "example.com"},
		Backend: &http.Client{Transport: transport},
	})
	client := NewClient(cfg)

	var parsed payload

	err := client.DoHTTP(http.MethodGet, "/", nil, &parsed)

	assert.Nil(t, err)
	assert.Equal(t, 5, parsed.Foo)
}
