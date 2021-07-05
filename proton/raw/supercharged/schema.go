package supercharged

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Request formalizes the fields for a canonical Supercharged request.
type Request struct {
	BaseURL  *url.URL
	Client   string
	Version  string
	Method   string
	Endpoint string
	Data     interface{}
}

// HTTPRequest creates an *http.Request from a Supercharged request descriptor.
func (r *Request) HTTPRequest() (*http.Request, error) {
	var req *http.Request

	if r.BaseURL == nil {
		return nil, errors.New("supercharged: missing request base URL")
	}

	data := r.Data
	if data == nil {
		data = struct{}{}
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resource, err := r.BaseURL.Parse(r.Endpoint)
	if err != nil {
		return nil, err
	}

	// By Supercharged conventions, the request data payload is stored in a request header for
	// GET and HEAD requests, and in the request body for other methods.
	switch r.Method {
	case http.MethodHead, http.MethodGet:
		req, err = http.NewRequest(r.Method, resource.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add(requestHeaderClientID, r.Client)
		req.Header.Add(requestHeaderClientVersion, r.Version)
		req.Header.Add(requestHeaderData, string(reqBody))
		req.Header.Add("User-Agent", fmt.Sprintf(clientUserAgent, r.Client, r.Version))
		return req, nil

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		req, err = http.NewRequest(r.Method, resource.String(), bytes.NewReader(reqBody))
		if err != nil {
			return nil, err
		}
		req.Header.Add(requestHeaderClientID, r.Client)
		req.Header.Add(requestHeaderClientVersion, r.Version)
		req.Header.Add("User-Agent", fmt.Sprintf(clientUserAgent, r.Client, r.Version))
		return req, nil

	default:
		return nil, errors.New("supercharged: unsupported HTTP request method")
	}
}

// WebsocketURL creates a *url.URL target suitable for websocket connections from a Supercharged
// request descriptor.
func (r *Request) WebsocketURL() (*url.URL, error) {
	var payloadData string

	if r.BaseURL == nil {
		return nil, errors.New("supercharged: missing request base URL")
	}

	if r.Data != nil {
		serializedData, err := json.Marshal(r.Data)
		if err != nil {
			return nil, err
		}

		payloadData = base64.StdEncoding.EncodeToString(serializedData)
	}

	wsURL := r.BaseURL
	wsURL.Path = r.Endpoint

	switch r.BaseURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	default:
		return nil, errors.New("supercharged: unknown URL scheme for websocket connection")
	}

	query := wsURL.Query()
	query.Set("client", r.Client)
	query.Set("version", r.Version)
	if payloadData != "" {
		query.Set("data", payloadData)
	}
	wsURL.RawQuery = query.Encode()

	return wsURL, nil
}

// Response formalizes a canonical Supercharged JSON response body.
type Response struct {
	Success bool            `json:"success"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// ParseResponse is a convenience function for parsing a Supercharged JSON-serialized body into a
// Response struct.
func ParseResponse(data []byte) (*Response, error) {
	var resp Response

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
