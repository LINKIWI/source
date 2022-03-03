package backend

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Threshold is a Backend that enforces a minimum threshold size for streaming puts and a maximum
// threshold size for unary puts.
//
// Streaming puts that fail to meet the minimum size threshold are instead converted into unary
// writes and flushed with the backend's PutObject API. Streaming puts that do satisfy the minimum
// size threshold are transparently passed to the underlying backend. Note that usage of this
// backend implies a maximum memory usage equal to the threshold size, since all pending (not yet
// flushed) writes are buffered in memory.
//
// Unary puts that exceed the maximum size threshold are similarly converted into streaming writes
// with the backend's StreamPutObject API. Unary puts that do satisfy the maximum size threshold are
// transparently passed to the underlying backend.
type Threshold struct {
	minStreamPutSize int
	maxUnaryPutSize  int
	Backend
}

// NewThreshold creates a Threshold wrapping an existing Backend with the desired threshold sizes,
// in bytes.
func NewThreshold(minStreamPutSize int, maxUnaryPutSize int, backend Backend) Backend {
	return &Threshold{
		minStreamPutSize: minStreamPutSize,
		maxUnaryPutSize:  maxUnaryPutSize,
		Backend:          backend,
	}
}

// PutObject transparently converts the unary put request to a streaming put request chunked at
// maxUnaryPutSize bytes if the total request data size exceeds the maxUnaryPutSize limit.
func (t *Threshold) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	stream := make(chan *service.PutObjectStreamRequest)

	if t.maxUnaryPutSize == 0 || len(request.Data) < t.maxUnaryPutSize {
		return t.Backend.PutObject(ctx, request)
	}

	go func() {
		var offset int

		// Clone the original request for asynchronous relay to the put stream, but remove
		// the attached data payload in order to save memory copy overhead when
		// constructing individual stream request messages.
		proxyRequest := proto.Clone(request).(*service.PutObjectRequest)
		proxyRequest.Data = nil

		for {
			if offset >= len(request.Data) {
				break
			}

			end := offset + t.maxUnaryPutSize
			if end > len(request.Data) {
				end = len(request.Data)
			}

			message := proto.Clone(proxyRequest).(*service.PutObjectRequest)
			message.Data = request.Data[offset:end]
			stream <- &service.PutObjectStreamRequest{Request: message}

			offset += t.maxUnaryPutSize
		}

		close(stream)
	}()

	response, err := t.Backend.StreamPutObject(ctx, stream)
	if err != nil {
		return nil, err
	}

	return response.Response, nil
}

// StreamPutObject buffers up to a maximum of threshold bytes in memory, followed by either draining
// the entire buffer as a single request to the backend's PutObject API, or passing all the chunks
// without modification to the underlying backend.
func (t *Threshold) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	var size int
	var buf []*service.PutObjectStreamRequest

	for request := range stream {
		buf = append(buf, request)
		size += len(request.Request.Data)

		// Don't consume any additional requests if the threshold has already been reached.
		if size >= t.minStreamPutSize {
			break
		}
	}

	// The total size of all chunks does not reach the threshold; combine the queue of buffered
	// chunk writes in order into a single request passed to the backend's non-streaming put.
	if size < t.minStreamPutSize {
		var request *service.PutObjectStreamRequest

		combined := make([]byte, 0, size)

		for _, request = range buf {
			combined = append(combined, request.Request.Data...)
		}

		unary := proto.Clone(request.Request).(*service.PutObjectRequest)
		unary.Data = combined

		resp, err := t.Backend.PutObject(ctx, unary)
		if err != nil {
			return nil, err
		}

		return &service.PutObjectStreamResponse{Response: resp}, nil
	}

	// The desired streaming size threshold has been reached; drain the queue of
	// buffered chunk writes in order, followed by redirecting any remaining chunks in
	// the stream to the proxy transparently.
	proxy := make(chan *service.PutObjectStreamRequest)

	go func() {
		for _, request := range buf {
			proxy <- request
		}

		for request := range stream {
			proxy <- request
		}

		close(proxy)
	}()

	return t.Backend.StreamPutObject(ctx, proxy)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (t *Threshold) Descriptor() *common.Node {
	return &common.Node{
		Name: "threshold",
		Children: map[string]*common.Node_Value{
			"min_stream_put_size": {
				Child: &common.Node_Value_Value{
					Value: fmt.Sprintf("%d", t.minStreamPutSize),
				},
			},
			"max_unary_put_size": {
				Child: &common.Node_Value_Value{
					Value: fmt.Sprintf("%d", t.maxUnaryPutSize),
				},
			},
			"backend": {
				Child: &common.Node_Value_Node{
					Node: t.Backend.Descriptor(),
				},
			},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (t *Threshold) String() string {
	return schemas.MarshalDescriptor(t.Descriptor())
}
