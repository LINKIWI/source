package webgrep

import (
	"net/http"

	"lib.kevinlin.info/proton"
)

// Client is a webgrep API client; effectively, a single layer of abstraction above a Supercharged
// HTTP client provided by Proton.
type Client struct {
	sc *proton.Client
}

// NewClient creates a new webgrep API client for an instance hosted at a particular base URL.
func NewClient(cfg *proton.Config) (*Client, error) {
	parsed, err := proton.NewConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{proton.NewClient(parsed)}, nil
}

// Search executes a search query.
func (c *Client) Search(request *SearchQueryRequest) (*SearchQueryResponse, error) {
	var resp SearchQueryResponse

	if err := c.sc.DoHTTP(http.MethodGet, EndpointSearch, request, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Metadata requests metadata about the webgrep instance.
func (c *Client) Metadata() (*MetadataResponse, error) {
	var resp MetadataResponse

	if err := c.sc.DoHTTP(http.MethodGet, EndpointMetadata, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
