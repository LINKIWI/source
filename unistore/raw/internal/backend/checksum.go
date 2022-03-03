package backend

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Checksum is a Backend that computes checksums for outgoing get object requests and validates
// client-provided checksums for incoming put object requests. In general, this middleware should be
// applied at the highest possible layer so that the checksum is most likely to represent the data
// the client interacts with (i.e. without modification by other middleware layers in the chain).
type Checksum struct {
	algorithm common.Checksum
	Backend
}

// NewChecksum creates a Checksum with the specified algorithm.
func NewChecksum(algorithm string, backend Backend) Backend {
	c := &Checksum{
		algorithm: common.Checksum(common.Checksum_value[strings.ToUpper(algorithm)]),
		Backend:   backend,
	}

	return newInstrumentation("checksum", c)
}

// GetObject defers to the underlying Backend and transparently injects a checksum in the response
// metadata.
func (c *Checksum) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	resp, err := c.Backend.GetObject(ctx, request)
	if err != nil {
		return nil, err
	}

	resp.Metadata.Checksum = c.checksum(resp.Data)

	return resp, nil
}

// StreamGetObject defers to the underlying Backend and transparently injects a checksum for each
// individual response in the stream.
func (c *Checksum) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	proxy := make(chan *service.GetObjectStreamResponse)

	stream, errs := c.Backend.StreamGetObject(ctx, request)
	if stream == nil {
		return nil, errs
	}

	go func() {
		for response := range stream {
			response.Response.Metadata.Checksum = c.checksum(response.Response.Data)

			proxy <- response
		}

		close(proxy)
	}()

	return proxy, errs
}

// PutObject validates the checksum in the request, if available. If the checksum matches, the
// request is passed to the underlying backend; otherwise, an error is returned.
func (c *Checksum) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	if request.Checksum != "" {
		actual := c.checksum(request.Data)

		if actual != request.Checksum {
			return nil, status.Errorf(
				codes.FailedPrecondition,
				"backend: data checksum mismatch: actual=%s expected=%s",
				actual,
				request.Checksum,
			)
		}
	}

	return c.Backend.PutObject(ctx, request)
}

// StreamPutObject validates the checksum for each individual request as it arrives over the stream,
// followed by passing the request to the underlying backend. Checksum validation failures result in
// early termination of the stream.
func (c *Checksum) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)
	responses := make(chan *service.PutObjectStreamResponse)
	errs := make(chan error)

	go func() {
		for request := range stream {
			if request.Request.Checksum != "" {
				actual := c.checksum(request.Request.Data)

				if actual != request.Request.Checksum {
					errs <- status.Errorf(
						codes.FailedPrecondition,
						"backend: data checksum mismatch: actual=%s expected=%s",
						actual,
						request.Request.Checksum,
					)
					close(proxy)
					return
				}
			}

			proxy <- request
		}

		close(proxy)
	}()

	go func() {
		response, err := c.Backend.StreamPutObject(ctx, proxy)
		if err != nil {
			errs <- err
		} else {
			responses <- response
		}
	}()

	select {
	case err := <-errs:
		return nil, err
	case response := <-responses:
		return response, nil
	}
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (c *Checksum) Descriptor() *common.Node {
	return &common.Node{
		Name: "checksum",
		Children: map[string]*common.Node_Value{
			"algorithm": {
				Child: &common.Node_Value_Value{
					Value: c.algorithm.String(),
				},
			},
			"backend": {
				Child: &common.Node_Value_Node{
					Node: c.Backend.Descriptor(),
				},
			},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (c *Checksum) String() string {
	return schemas.MarshalDescriptor(c.Descriptor())
}

// checksum computes a checksum for the provided data slice.
func (c *Checksum) checksum(data []byte) string {
	var h hash.Hash

	switch c.algorithm {
	case common.Checksum_CRC32:
		h = crc32.NewIEEE()
	case common.Checksum_MD5:
		h = md5.New()
	case common.Checksum_SHA1:
		h = sha1.New()
	case common.Checksum_SHA256:
		h = sha256.New()
	case common.Checksum_SHA512:
		h = sha512.New()
	default:
		return ""
	}

	h.Write(data)

	return fmt.Sprintf("%x", h.Sum(nil))
}
