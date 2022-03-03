package b2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// Option aliases a callback function to apply options on the client.
type Option func(*Client)

// Client is a structured wrapper for the B2 API. Note that it provides minimal additional features
// and mostly intends to provide a 1:1 mapping to B2 HTTP APIs.
type Client struct {
	http     *http.Client
	log      *zap.Logger
	identity string
}

// New creates a new B2 client with an optional set of options.
func New(opts ...Option) (*Client, error) {
	client := &Client{
		http:     http.DefaultClient,
		log:      zap.NewNop(),
		identity: DefaultClientID,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// WithTransport provides a custom HTTP client transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.http.Transport = transport
	}
}

// WithLogger enables verbose logging to the specified logger.
func WithLogger(log *zap.Logger) Option {
	return func(c *Client) {
		c.log = log
	}
}

// WithClientID specifies a client identity to use as the User-Agent for outbound HTTP requests to
// the B2 API.
func WithClientID(id string) Option {
	return func(c *Client) {
		c.identity = id
	}
}

// AuthorizeAccount implements the b2_authorize_account API.
// Reference: https://www.backblaze.com/b2/docs/b2_authorize_account.html
func (c *Client) AuthorizeAccount(ctx context.Context, request *AuthorizeAccountRequest) (*AuthorizeAccountResponse, error) {
	var response AuthorizeAccountResponse

	u, err := url.Parse(APIV2BaseURL)
	if err != nil {
		return nil, statusFromError(err)
	}

	r, err := c.createAPIRequest(ctx, u, EndpointAuthorizeAccount, http.MethodGet, request)
	if err != nil {
		return nil, err
	}

	r.SetBasicAuth(request.ApplicationKeyID, request.ApplicationKey)

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelLargeFile implements the b2_cancel_large_file API.
// Reference: https://www.backblaze.com/b2/docs/b2_cancel_large_file.html
func (c *Client) CancelLargeFile(ctx context.Context, request *CancelLargeFileRequest) (*CancelLargeFileResponse, error) {
	var response CancelLargeFileResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointCancelLargeFile, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CopyFile implements the b2_copy_file API.
// Reference: https://www.backblaze.com/b2/docs/b2_copy_file.html
func (c *Client) CopyFile(ctx context.Context, request *CopyFileRequest) (*CopyFileResponse, error) {
	var response CopyFileResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointCopyFile, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CopyPart implements the b2_copy_part API.
// Reference: https://www.backblaze.com/b2/docs/b2_copy_part.html
func (c *Client) CopyPart(ctx context.Context, request *CopyPartRequest) (*CopyPartResponse, error) {
	var response CopyPartResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointCopyPart, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteFileVersion implements the b2_delete_file_version API.
// Reference: https://www.backblaze.com/b2/docs/b2_delete_file_version.html
func (c *Client) DeleteFileVersion(ctx context.Context, request *DeleteFileVersionRequest) (*DeleteFileVersionResponse, error) {
	var response DeleteFileVersionResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointDeleteFileVersion, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DownloadFileByName implements the b2_download_file_by_name API.
// Reference: https://www.backblaze.com/b2/docs/b2_download_file_by_name.html
func (c *Client) DownloadFileByName(ctx context.Context, request *DownloadFileByNameRequest) (*DownloadFileByNameResponse, error) {
	r, err := c.createAPIRequest(ctx, request.DownloadURL, "", http.MethodGet, request)
	if err != nil {
		return nil, err
	}

	ref := &url.URL{Path: path.Join("file", request.Bucket, request.File)}
	r.URL = r.URL.ResolveReference(ref)

	if request.Range.Start > 0 && request.Range.End > 0 {
		r.Header.Set(
			"Range",
			fmt.Sprintf("%s=%d-%d", request.Range.Unit, request.Range.Start, request.Range.End),
		)
	} else if request.Range.Start > 0 {
		r.Header.Set(
			"Range",
			fmt.Sprintf("%s=%d-", request.Range.Unit, request.Range.Start),
		)
	} else if request.Range.End > 0 {
		r.Header.Set(
			"Range",
			fmt.Sprintf("%s=0-%d", request.Range.Unit, request.Range.End),
		)
	}

	// Note that the response body is not closed as this layer since it is passed directly to
	// the client as an io.ReadCloser; it is the responsibility of the caller to close the
	// body when it has been consumed.
	response, err := c.doRaw(r)
	if err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		return nil, statusFromError(err)
	}

	contentRange := Range{}
	contentRangeRe := regexp.MustCompile("^(?P<unit>[^\\s]+)\\s(?P<start>\\d+)-(?P<end>\\d+)/(?P<total>\\d+)$")
	matches := contentRangeRe.FindStringSubmatch(response.Header.Get("Content-Range"))
	if len(matches) == 5 {
		start, err := strconv.Atoi(matches[contentRangeRe.SubexpIndex("start")])
		if err != nil {
			return nil, statusFromError(err)
		}

		end, err := strconv.Atoi(matches[contentRangeRe.SubexpIndex("end")])
		if err != nil {
			return nil, statusFromError(err)
		}

		total, err := strconv.Atoi(matches[contentRangeRe.SubexpIndex("total")])
		if err != nil {
			return nil, statusFromError(err)
		}

		contentRange.Unit = matches[contentRangeRe.SubexpIndex("unit")]
		contentRange.Start = start
		contentRange.End = end
		contentRange.Total = total
	}

	fileInfo := make(map[string]string)
	for header := range response.Header {
		if strings.HasPrefix(header, "X-Bz-Info-") {
			fileInfo[strings.ToLower(header[len("X-Bz-Info-"):])] = response.Header.Get(header)
		}
	}

	var uploadTimestamp UnixTimestamp
	if err := json.Unmarshal([]byte(response.Header.Get("X-Bz-Upload-Timestamp")), &uploadTimestamp); err != nil {
		return nil, statusFromError(err)
	}

	return &DownloadFileByNameResponse{
		FileID:          response.Header.Get("X-Bz-File-Id"),
		FileName:        response.Header.Get("X-Bz-File-Name"),
		FileInfo:        fileInfo,
		ContentSHA1:     response.Header.Get("X-Bz-Content-Sha1"),
		Size:            size,
		ContentType:     response.Header.Get("Content-Type"),
		Range:           contentRange,
		UploadTimestamp: &uploadTimestamp,
		Data:            response.Body,
	}, nil
}

// FinishLargeFile implements the b2_finish_large_file API.
// Reference: https://www.backblaze.com/b2/docs/b2_finish_large_file.html
func (c *Client) FinishLargeFile(ctx context.Context, request *FinishLargeFileRequest) (*FinishLargeFileResponse, error) {
	var response FinishLargeFileResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointFinishLargeFile, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUploadPartURL implements the b2_get_upload_part_url API.
// Reference: https://www.backblaze.com/b2/docs/b2_get_upload_part_url.html
func (c *Client) GetUploadPartURL(ctx context.Context, request *GetUploadPartURLRequest) (*GetUploadPartURLResponse, error) {
	var response GetUploadPartURLResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointGetUploadPartURL, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUploadURL implements the b2_get_upload_url API.
// Reference: https://www.backblaze.com/b2/docs/b2_get_upload_url.html
func (c *Client) GetUploadURL(ctx context.Context, request *GetUploadURLRequest) (*GetUploadURLResponse, error) {
	var response GetUploadURLResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointGetUploadURL, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListBuckets implements the b2_list_buckets API.
// Reference: https://www.backblaze.com/b2/docs/b2_list_buckets.html
func (c *Client) ListBuckets(ctx context.Context, request *ListBucketsRequest) (*ListBucketsResponse, error) {
	var response ListBucketsResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointListBuckets, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListFileNames implements the b2_list_file_names API.
// Reference: https://www.backblaze.com/b2/docs/b2_list_file_names.html
func (c *Client) ListFileNames(ctx context.Context, request *ListFileNamesRequest) (*ListFileNamesResponse, error) {
	var response ListFileNamesResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointListFileNames, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// StartLargeFile implements the b2_start_large_file API.
// Reference: https://www.backblaze.com/b2/docs/b2_start_large_file.html
func (c *Client) StartLargeFile(ctx context.Context, request *StartLargeFileRequest) (*StartLargeFileResponse, error) {
	var response StartLargeFileResponse

	r, err := c.createAPIRequest(ctx, request.BaseURL, EndpointStartLargeFile, http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UploadFile implements the b2_upload_file API.
// Reference: https://www.backblaze.com/b2/docs/b2_upload_file.html
func (c *Client) UploadFile(ctx context.Context, request *UploadFileRequest) (*UploadFileResponse, error) {
	var response UploadFileResponse

	r, err := c.createAPIRequest(ctx, request.UploadURL, "", http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	if request.FileInfo != nil {
		for key, info := range request.FileInfo {
			r.Header.Set(fmt.Sprintf("X-Bz-Info-%s", key), url.QueryEscape(info))
		}
	}

	// Specify the content length explicitly to prevent net/http from automatically converting
	// the request to use chunked transfer encoding, which is not supported by the B2 API for
	// uploads.
	r.ContentLength = int64(request.ContentLength)
	r.Body = request.Data

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, err
}

// UploadPart implements the b2_upload_part API.
// Reference: https://www.backblaze.com/b2/docs/b2_upload_part.html
func (c *Client) UploadPart(ctx context.Context, request *UploadPartRequest) (*UploadPartResponse, error) {
	var response UploadPartResponse

	r, err := c.createAPIRequest(ctx, request.UploadURL, "", http.MethodPost, request)
	if err != nil {
		return nil, err
	}

	r.ContentLength = int64(request.ContentLength)
	r.Body = request.Data

	if _, err := c.doAPI(r, &response); err != nil {
		return nil, err
	}

	return &response, err
}

// createAPIRequest creates a http.Request object, obeying B2 conventions, based on the specified
// base URL, endpoint, HTTP method, and request struct.
func (c *Client) createAPIRequest(ctx context.Context, baseURL *url.URL, endpoint string, method string, request interface{}) (*http.Request, error) {
	var body io.ReadCloser

	u := baseURL
	if endpoint != "" {
		u = baseURL.ResolveReference(&url.URL{Path: path.Join(APIV2BasePath, endpoint)})
	}

	headers, err := marshalHeaders(request)
	if err != nil {
		return nil, statusFromError(err)
	}

	headers.Set("User-Agent", c.identity)

	if method == http.MethodPost {
		payload, err := json.Marshal(request)
		if err != nil {
			return nil, statusFromError(err)
		}

		body = io.NopCloser(bytes.NewReader(payload))
	}

	r := &http.Request{
		Method: method,
		URL:    u,
		Header: headers,
		Body:   body,
	}

	return r.WithContext(ctx), nil
}

// doAPI executes the request as a canonical API request, whose response body can be parsed as JSON
// into the provided container.
func (c *Client) doAPI(request *http.Request, data interface{}) (*http.Response, error) {
	response, err := c.doRaw(request)
	if err != nil {
		return nil, statusFromError(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return response, statusFromError(err)
	}

	if data != nil {
		if err := json.Unmarshal(body, &data); err != nil {
			return response, statusFromError(err)
		}
	}

	return response, nil
}

// doRaw executes the request, handling HTTP errors semantically as API errors, and returns the
// response as-is without any response body manipulation.
func (c *Client) doRaw(request *http.Request) (*http.Response, error) {
	var status Status

	c.log.Debug(
		"dispatching HTTP request",
		zap.String("method", request.Method),
		zap.Any("headers", request.Header),
		zap.Stringer("url", request.URL),
	)

	response, err := c.http.Do(request)
	if err != nil {
		c.log.Error(
			"fatal transport error during HTTP dispatch",
			zap.String("method", request.Method),
			zap.Stringer("url", request.URL),
			zap.Error(err),
		)

		return nil, statusFromError(err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return response, statusFromError(err)
		}

		if err := json.Unmarshal(body, &status); err != nil {
			return response, statusFromError(err)
		}

		response.Body.Close()

		c.log.Debug(
			"received application-level error from B2 API",
			zap.Int("status", status.Status),
			zap.String("code", status.Code),
			zap.String("message", status.Message),
		)

		return response, &status
	}

	c.log.Debug(
		"successful B2 API response",
		zap.String("method", response.Request.Method),
		zap.Stringer("url", response.Request.URL),
		zap.Int("status", response.StatusCode),
		zap.Int64("size", response.ContentLength),
	)

	return response, nil
}
