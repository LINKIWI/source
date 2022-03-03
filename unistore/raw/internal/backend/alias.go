package backend

import (
	"context"

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
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.HeadBucket(ctx, request)
}

// HeadObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.HeadObject(ctx, request)
}

// GetObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.GetObject(ctx, request)
}

// StreamGetObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Request.Resource.Bucket) {
			request.Request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.StreamGetObject(ctx, request)
}

// PutObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.PutObject(ctx, request)
}

// StreamPutObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)

	go func() {
		for request := range stream {
			for _, alias := range a.aliases {
				if alias.Pattern.MatchString(request.Request.Resource.Bucket) {
					request.Request.Resource.Bucket = alias.Alias
					break
				}
			}

			proxy <- request
		}

		close(proxy)
	}()

	return a.Backend.StreamPutObject(ctx, proxy)
}

// DeleteObject conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.DeleteObject(ctx, request)
}

// ListObjects conditionally overwrites the request bucket and defers to the underlying backend.
func (a *Alias) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	for _, alias := range a.aliases {
		if alias.Pattern.MatchString(request.Resource.Bucket) {
			request.Resource.Bucket = alias.Alias
			break
		}
	}

	return a.Backend.ListObjects(ctx, request)
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
