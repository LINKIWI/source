package backend

import (
	"context"

	"google.golang.org/protobuf/proto"

	"unistore/internal/config"
	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Alias is a Backend that optionally replaces bucket names according to an aliases definition that
// maps regular expression patterns to their canonical bucket names.
type Alias struct {
	aliases []*config.Alias
	Backend
}

// NewAlias creates a new Alias with the specified canonical name mapping.
func NewAlias(aliases []*config.Alias, backend Backend) Backend {
	return &Alias{
		aliases: aliases,
		Backend: backend,
	}
}

// HeadBucket conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.HeadBucketRequest)
			proxyRequest.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.HeadBucket(ctx, proxyRequest)
}

// HeadObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.HeadObjectRequest)
			proxyRequest.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.HeadObject(ctx, proxyRequest)
}

// GetObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.GetObjectRequest)
			proxyRequest.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.GetObject(ctx, proxyRequest)
}

// StreamGetObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.GetObjectStreamRequest)
			proxyRequest.Request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.StreamGetObject(ctx, proxyRequest)
}

// PutObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.PutObjectRequest)
			proxyRequest.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.PutObject(ctx, proxyRequest)
}

// StreamPutObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)

	go func() {
		for request := range stream {
			proxyRequest := request

			for _, alias := range a.aliases {
				if alias.Pattern.MatchString(request.Request.Resource.Bucket) {
					proxyRequest = proto.Clone(request).(*service.PutObjectStreamRequest)
					proxyRequest.Request.Resource.Bucket = alias.Alias
					break
				}
			}

			proxy <- proxyRequest
		}

		close(proxy)
	}()

	return a.Backend.StreamPutObject(ctx, proxy)
}

// DeleteObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.DeleteObjectRequest)
			proxyRequest.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.DeleteObject(ctx, proxyRequest)
}

// ListObjects conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	proxyRequest := request

	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			proxyRequest = proto.Clone(request).(*service.ListObjectsRequest)
			proxyRequest.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.ListObjects(ctx, proxyRequest)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (a *Alias) Descriptor() *common.Node {
	return &common.Node{
		Name: "alias",
		Children: map[string]*common.Node_Value{
			"backend": {Child: &common.Node_Value_Node{Node: a.Backend.Descriptor()}},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (a *Alias) String() string {
	return schemas.MarshalDescriptor(a.Descriptor())
}
