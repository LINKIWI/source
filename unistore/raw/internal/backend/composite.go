package backend

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Composite is a Backend that abstracts over one or more children backends with configurable
// dispatch policies partitioned along read and write RPCs. It is useful for high-availability
// deployments where data redundancy and/or online failover is required.
type Composite struct {
	readDispatch  common.Dispatch
	writeDispatch common.Dispatch
	backends      []Backend
}

// NewComposite creates a new Composite with the specified dispatch policies and children.
func NewComposite(readDispatch string, writeDispatch string, backends []Backend) Backend {
	c := &Composite{
		readDispatch:  common.Dispatch(common.Dispatch_value[strings.ToUpper(readDispatch)]),
		writeDispatch: common.Dispatch(common.Dispatch_value[strings.ToUpper(writeDispatch)]),
		backends:      backends,
	}

	return newInstrumentation("composite", c)
}

// HeadBucket invokes HeadBucket in one or more composed backends.
func (c *Composite) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.HeadBucket(ctx, proto.Clone(request).(*service.HeadBucketRequest))
	}

	response, err := c.dispatch(c.readDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.HeadBucketResponse), nil
}

// HeadObject invokes HeadObject in one or more composed backends.
func (c *Composite) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.HeadObject(ctx, proto.Clone(request).(*service.HeadObjectRequest))
	}

	response, err := c.dispatch(c.readDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.HeadObjectResponse), nil
}

// GetObject invokes GetObject in one or more composed backends.
func (c *Composite) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.GetObject(ctx, proto.Clone(request).(*service.GetObjectRequest))
	}

	response, err := c.dispatch(c.readDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.GetObjectResponse), nil
}

// StreamGetObject invokes StreamGetObject in one or more composed backends.
// Note that only specific dispatch policies are implemented, and composite streaming gets come with
// a set of unique behavior nuances; see the inline commentary for details.
func (c *Composite) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	// StreamGetObject is a bit of a special case. In practice, only one backend serves the
	// response that is eventually consumed by the client; dispatching the request to multiple
	// backends would result in a stuck coroutine (due to an absence of channel consumption).
	// Thus, only a few dispatch policies are supported by this RPC.

	responseProxy := make(chan *service.GetObjectStreamResponse)
	errsProxy := make(chan error, 1)

	switch c.readDispatch {
	case common.Dispatch_RANDOM:
		return c.backends[rand.Intn(len(c.backends))].StreamGetObject(ctx, request)

	case common.Dispatch_INITIAL:
		return c.backends[0].StreamGetObject(ctx, request)

	case common.Dispatch_FAILOVER:
		go func() {
			for _, backend := range c.backends {
				responses, errs := backend.StreamGetObject(ctx, request)
				if responses == nil {
					// It must be the case that there is an error; don't attempt
					// to consume responses and just try the next backend.
					continue
				}

				// At this point, the stream is non-nil and any potential future
				// errors would be encountered mid-stream. However, since data is
				// already being exchanged at this point, failover is no longer
				// permissible as it would semantically equivalent to restarting the
				// data stream. Instead, propagate any encountered errors directly
				// to the client (which effectively errors out the entire request).

				select {
				case err := <-errs:
					errsProxy <- err
					return
				case response, ok := <-responses:
					if !ok {
						close(responseProxy)
						return
					}

					responseProxy <- response
				}
			}
		}()

		return responseProxy, errsProxy

	default:
		errsProxy <- status.Errorf(
			codes.Unimplemented,
			"backend: unsupported dispatch policy for streaming get: policy=%v",
			c.readDispatch,
		)

		return nil, errsProxy
	}
}

// PutObject invokes PutObject in one or more composed backends.
func (c *Composite) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.PutObject(ctx, proto.Clone(request).(*service.PutObjectRequest))
	}

	response, err := c.dispatch(c.writeDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.PutObjectResponse), nil
}

// StreamPutObject invokes StreamPutObject in one or more composed backends.
func (c *Composite) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	// Synchronous dispatch is inherently incompatible with streaming writes. Semantically, one
	// would expect for an ALL_SYNC policy to complete the full stream of writes to the nth
	// backend, followed by repeating the same sequence of writes to subsequent (n + 1)th
	// backends. However, this is impossible to accomplish unless the server buffered the full
	// sequence of chunks from the entire stream in memory, which defeats the purpose of a
	// streaming API and runs counter to its coroutine-oriented architecture. Thus, simply
	// reject this dispatch policy altogether.
	if c.writeDispatch == common.Dispatch_ALL_SYNC {
		return nil, status.Error(
			codes.Unimplemented,
			"backend: ALL_SYNC dispatch policy is not supported for streaming puts",
		)
	}

	// Duplicate the inbound request to separate channels, one per backend.
	proxy := make(map[Backend]chan *service.PutObjectStreamRequest)
	for _, backend := range c.backends {
		proxy[backend] = make(chan *service.PutObjectStreamRequest)
	}

	go func() {
		for request := range stream {
			for _, backend := range c.backends {
				proxy[backend] <- proto.Clone(request).(*service.PutObjectStreamRequest)
			}
		}

		for _, backend := range c.backends {
			close(proxy[backend])
		}
	}()

	rpc := func(backend Backend) (proto.Message, error) {
		return backend.StreamPutObject(ctx, proxy[backend])
	}

	response, err := c.dispatch(c.writeDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.PutObjectStreamResponse), nil
}

// DeleteObject invokes DeleteObject in one or more composed backends.
func (c *Composite) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.DeleteObject(ctx, proto.Clone(request).(*service.DeleteObjectRequest))
	}

	response, err := c.dispatch(c.writeDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.DeleteObjectResponse), nil
}

// ListBuckets invokes ListBuckets in one or more composed backends.
func (c *Composite) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.ListBuckets(ctx, proto.Clone(request).(*service.ListBucketsRequest))
	}

	response, err := c.dispatch(c.readDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.ListBucketsResponse), nil
}

// ListObjects invokes ListObjects in one or more composed backends.
func (c *Composite) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	rpc := func(backend Backend) (proto.Message, error) {
		return backend.ListObjects(ctx, proto.Clone(request).(*service.ListObjectsRequest))
	}

	response, err := c.dispatch(c.readDispatch, rpc)
	if err != nil {
		return nil, err
	}

	return response.(*service.ListObjectsResponse), nil
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (c *Composite) Descriptor() *common.Node {
	children := map[string]*common.Node_Value{
		"read_dispatch": {
			Child: &common.Node_Value_Value{
				Value: c.readDispatch.String(),
			},
		},
		"write_dispatch": {
			Child: &common.Node_Value_Value{
				Value: c.writeDispatch.String(),
			},
		},
	}

	for i, backend := range c.backends {
		children[fmt.Sprintf("backend_%d", i)] = &common.Node_Value{
			Child: &common.Node_Value_Node{
				Node: backend.Descriptor(),
			},
		}
	}

	return &common.Node{
		Name:     "composite",
		Children: children,
	}
}

// String returns a human-consumable representation of this backend.
func (c *Composite) String() string {
	return schemas.MarshalDescriptor(c.Descriptor())
}

// dispatch is a convenience method to dispatch a request based on the selected dispatch policy.
func (c *Composite) dispatch(policy common.Dispatch, rpc func(backend Backend) (proto.Message, error)) (proto.Message, error) {
	switch policy {
	case common.Dispatch_RANDOM:
		return c.dispatchRandom(rpc)
	case common.Dispatch_INITIAL:
		return c.dispatchInitial(rpc)
	case common.Dispatch_FAILOVER:
		return c.dispatchFailover(rpc)
	case common.Dispatch_ALL_SYNC:
		return c.dispatchAllSync(rpc)
	case common.Dispatch_ALL_ASYNC:
		return c.dispatchAllAsync(rpc)
	default:
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: unsupported dispatch policy for method: policy=%v",
			policy,
		)
	}
}

// dispatchRandom is a dispatcher implementation that dispatches the request to a randomly selected
// backend.
func (c *Composite) dispatchRandom(rpc func(backend Backend) (proto.Message, error)) (proto.Message, error) {
	return rpc(c.backends[rand.Intn(len(c.backends))])
}

// dispatchInitial is a dispatcher implementation that dispatches the request to the first backend.
func (c *Composite) dispatchInitial(rpc func(backend Backend) (proto.Message, error)) (proto.Message, error) {
	return rpc(c.backends[0])
}

// dispatchFailover is a dispatcher implementation that attempts the request against backends
// sequentially until one provides a non-error result.
func (c *Composite) dispatchFailover(rpc func(backend Backend) (proto.Message, error)) (proto.Message, error) {
	var errs []error

	for _, backend := range c.backends {
		response, err := rpc(backend)
		if err == nil {
			return response, nil
		}

		zap.L().Warn(
			"error during invocation; advancing to next backend by failover policy",
			zap.Error(err),
		)

		errs = append(errs, err)
	}

	return nil, fmt.Errorf("backend: all backends errored in failover policy: errs=%v", errs)
}

// dispatchAllSync is a dispatcher implementation that synchronously dispatches the request to all
// backends, one after the other in sequential order.
func (c *Composite) dispatchAllSync(rpc func(backend Backend) (proto.Message, error)) (proto.Message, error) {
	var response proto.Message
	var err error

	for _, backend := range c.backends {
		response, err = rpc(backend)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// dispatchAllAsync is a dispatcher implementation that asynchronously dispatches the request to all
// backends and waits for all requests to complete. If any request errors, that error is propagated
// to the client. Otherwise, the first successful response is returned.
func (c *Composite) dispatchAllAsync(rpc func(backend Backend) (proto.Message, error)) (proto.Message, error) {
	responses := make(chan proto.Message, len(c.backends))
	errs := make(chan error, len(c.backends))

	for _, backend := range c.backends {
		go func(backend Backend) {
			response, err := rpc(backend)
			errs <- err

			if response != nil {
				responses <- response
			}
		}(backend)
	}

	// Block until all requests have been completed.
	for i := 0; i < len(c.backends); i++ {
		if err := <-errs; err != nil {
			return nil, err
		}
	}

	return <-responses, nil
}
