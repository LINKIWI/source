package backend

import (
	"context"
	"sort"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Mux is a Backend that demultiplexes multiple backends identified by a common.Backend enum.
type Mux struct {
	backends map[common.Backend]Backend
}

// NewMux creates a Mux from a backend implementation map.
func NewMux(backends map[common.Backend]Backend) Backend {
	return newInstrumentation("mux", &Mux{backends})
}

// HeadBucket demuxes HeadBucket.
func (m *Mux) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.HeadBucket(ctx, request)
}

// HeadObject demuxes HeadObject.
func (m *Mux) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.HeadObject(ctx, request)
}

// GetObject demuxes GetObject.
func (m *Mux) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.GetObject(ctx, request)
}

// StreamGetObject demuxes StreamGetObject.
func (m *Mux) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	errs := make(chan error, 1)

	backend, ok := m.backends[request.Request.Resource.Backend]
	if !ok {
		errs <- status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Request.Resource.Backend,
		)

		return nil, errs
	}

	return backend.StreamGetObject(ctx, request)
}

// PutObject demuxes PutObject.
func (m *Mux) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.PutObject(ctx, request)
}

// StreamPutObject demuxes StreamPutObject. Note that the backend metadata is included in the first
// stream request sent through the channel, so the implementation creates an intermediate "proxy"
// channel to inspect the correct backend to use based on the first request followed by passing it
// along to the backend implementation.
func (m *Mux) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)
	peek := <-stream

	if peek == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"backend: error reading first chunk from stream",
		)
	}

	backend, ok := m.backends[peek.Request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			peek.Request.Resource.Backend,
		)
	}

	go func() {
		proxy <- peek
		for request := range stream {
			proxy <- request
		}
		close(proxy)
	}()

	return backend.StreamPutObject(ctx, proxy)
}

// DeleteObject demuxes DeleteObject.
func (m *Mux) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.DeleteObject(ctx, request)
}

// ListBuckets demuxes ListBuckets.
func (m *Mux) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.ListBuckets(ctx, request)
}

// ListObjects demuxes ListObjects.
func (m *Mux) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	backend, ok := m.backends[request.Resource.Backend]
	if !ok {
		return nil, status.Errorf(
			codes.Unimplemented,
			"backend: no implementation exists for the requested backend: backend=%v",
			request.Resource.Backend,
		)
	}

	return backend.ListObjects(ctx, request)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (m *Mux) Descriptor() *common.Node {
	ids := make([]string, 0, len(m.backends))
	children := make(map[string]*common.Node_Value, len(m.backends))

	for id := range m.backends {
		ids = append(ids, id.String())
	}
	sort.Strings(ids)

	for id, backend := range m.backends {
		children[id.String()] = &common.Node_Value{
			Child: &common.Node_Value_Node{
				Node: backend.Descriptor(),
			},
		}
	}

	children["backends"] = &common.Node_Value{
		Child: &common.Node_Value_Value{
			Value: strings.Join(ids, ", "),
		},
	}

	return &common.Node{
		Name:     "mux",
		Children: children,
	}
}

// Close closes all backends asynchronously, waits for all closes to complete, and returns the first
// error if available, or nil if there are no errors.
func (m *Mux) Close() error {
	errs := make(chan error, len(m.backends))

	for _, backend := range m.backends {
		go func(backend Backend) {
			errs <- backend.Close()
		}(backend)
	}

	for i := 0; i < len(m.backends); i++ {
		if err := <-errs; err != nil {
			return err
		}
	}

	return nil
}

// String returns a human-consumable representation of this backend.
func (m *Mux) String() string {
	return schemas.MarshalDescriptor(m.Descriptor())
}
