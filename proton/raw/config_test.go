package proton

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig(&Config{
		BaseURL: &url.URL{Scheme: "https", Host: "example.com"},
	})

	assert.NoError(t, err)
	assert.Equal(t, "proton", cfg.ClientID)
	assert.Equal(t, "none", cfg.ClientVersion)
	assert.Equal(t, http.DefaultClient, cfg.Backend)
	assert.NotNil(t, cfg.Metrics)
	assert.NotNil(t, cfg.Logger)
}

func TestConfigValidationFailure(t *testing.T) {
	t.Parallel()

	cases := []*Config{
		// Missing base URL
		{},
		// Malformed base URL
		{
			BaseURL: &url.URL{},
		},
		{
			BaseURL: &url.URL{Scheme: "http"},
		},
		{
			BaseURL: &url.URL{Host: "example.com"},
		},
		// Client ID with illegal characters
		{
			BaseURL:  &url.URL{Scheme: "http", Host: "example.com"},
			ClientID: "#",
		},
		{
			BaseURL:  &url.URL{Scheme: "http", Host: "example.com"},
			ClientID: "client id",
		},
	}

	for _, cfg := range cases {
		assert.Error(t, cfg.validate())
	}
}

func TestConfigValidationSuccess(t *testing.T) {
	t.Parallel()

	cases := []*Config{
		{
			BaseURL: &url.URL{Scheme: "http", Host: "example.com"},
		},
		{
			BaseURL: &url.URL{Scheme: "http", Host: "example.com", Path: "/path"},
		},
	}

	for _, cfg := range cases {
		assert.NoError(t, cfg.validate())
	}
}
