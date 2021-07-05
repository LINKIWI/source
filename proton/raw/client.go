package proton

import (
	"context"
	"encoding/json"
	"io"

	"github.com/gorilla/websocket"
	"lib.kevinlin.info/aperture/lib"

	"lib.kevinlin.info/proton/supercharged"
)

// Client is a Supercharged client that transacts requests and responses over HTTP.
type Client struct {
	cfg *Config
}

// NewClient creates a new client instance for a Supercharged-compliant HTTP server, with an
// optional http.Client backend for initiating requests.
func NewClient(cfg *Config) *Client {
	return &Client{cfg}
}

// DoHTTP executes an HTTP request against the server.
func (c *Client) DoHTTP(method string, endpoint string, data interface{}, response interface{}) (scErr *supercharged.Error) {
	stopwatch := lib.NewStopwatch()
	tags := map[string]interface{}{
		"method":   method,
		"endpoint": endpoint,
		"client":   c.cfg.ClientID,
	}

	defer func() {
		c.cfg.Metrics.Incr("proton.http.request.invoke", tags)
		c.cfg.Metrics.Timing("proton.http.request.duration", stopwatch.Elapsed(), tags)

		if scErr != nil {
			c.cfg.Metrics.Incr("proton.http.request.error", tags)

			c.cfg.Logger.Printf(
				"[proton] transport or application HTTP error: "+
					"status=%d code=%d message=%s data=%+v",
				scErr.Status,
				scErr.Code,
				scErr.Message,
				scErr.Data,
			)
		}
	}()

	scReq := &supercharged.Request{
		BaseURL:  c.cfg.BaseURL,
		Client:   c.cfg.ClientID,
		Version:  c.cfg.ClientVersion,
		Method:   method,
		Endpoint: endpoint,
		Data:     data,
	}

	c.cfg.Logger.Printf(
		"[proton] sending HTTP request: "+
			"base_url=%v client=%s method=%s endpoint=%s data=%+v",
		scReq.BaseURL,
		scReq.Client,
		scReq.Method,
		scReq.Endpoint,
		scReq.Data,
	)

	httpReq, err := scReq.HTTPRequest()
	if err != nil {
		return supercharged.WrapError(err)
	}

	httpResp, err := c.cfg.Backend.Do(httpReq)
	if err != nil {
		return supercharged.WrapError(err)
	}
	tags["status"] = httpResp.StatusCode

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return supercharged.WrapError(err)
	}

	scResp, err := supercharged.ParseResponse(respBody)
	if err != nil {
		return supercharged.WrapError(err)
	}
	tags["code"] = scResp.Code

	if !scResp.Success {
		return &supercharged.Error{
			Status:  httpResp.StatusCode,
			Code:    scResp.Code,
			Message: scResp.Message,
			Data:    scResp.Data,
		}
	}

	if response != nil {
		if err := json.Unmarshal(scResp.Data, response); err != nil {
			return supercharged.WrapError(err)
		}
	}

	c.cfg.Metrics.Size("proton.http.request.body_size", httpReq.ContentLength, tags)
	c.cfg.Metrics.Size("proton.http.response.body_size", httpResp.ContentLength, tags)

	c.cfg.Logger.Printf(
		"[proton] successful HTTP response: status=%d code=%d message=%s data=%+v",
		httpResp.StatusCode,
		scResp.Code,
		scResp.Message,
		response,
	)

	return nil
}

// DoWebsocket establishes a duplex websocket connection with the server.
func (c *Client) DoWebsocket(endpoint string, data interface{}) (conn *websocket.Conn, scErr *supercharged.Error) {
	stopwatch := lib.NewStopwatch()
	tags := map[string]interface{}{
		"endpoint": endpoint,
		"client":   c.cfg.ClientID,
	}

	defer func() {
		c.cfg.Metrics.Incr("proton.websocket.connect.establish", tags)
		c.cfg.Metrics.Timing("proton.websocket.connect.duration", stopwatch.Elapsed(), tags)

		if scErr != nil {
			c.cfg.Metrics.Incr("proton.websocket.connect.error", tags)

			c.cfg.Logger.Printf(
				"[proton] error establishing websocket connection: "+
					"status=%d code=%d message=%s data=%+v",
				scErr.Status,
				scErr.Code,
				scErr.Message,
				scErr.Data,
			)
		}
	}()

	scReq := &supercharged.Request{
		BaseURL:  c.cfg.BaseURL,
		Client:   c.cfg.ClientID,
		Version:  c.cfg.ClientVersion,
		Endpoint: endpoint,
		Data:     data,
	}

	c.cfg.Logger.Printf(
		"[proton] dialing websocket connection: base_url=%v client=%s endpoint=%s data=%+v",
		scReq.BaseURL,
		scReq.Client,
		scReq.Endpoint,
		scReq.Data,
	)

	wsURL, err := scReq.WebsocketURL()
	if err != nil {
		return nil, supercharged.WrapError(err)
	}

	dialer := &websocket.Dialer{
		HandshakeTimeout:  c.cfg.Backend.Timeout,
		EnableCompression: true,
	}

	conn, _, err = dialer.DialContext(context.Background(), wsURL.String(), nil)
	if err != nil {
		return nil, supercharged.WrapError(err)
	}

	c.cfg.Logger.Printf(
		"[proton] successfully established websocket connection: local=%v remote=%v",
		conn.LocalAddr(),
		conn.RemoteAddr(),
	)

	return conn, nil
}
