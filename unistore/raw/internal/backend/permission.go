package backend

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"unistore/internal/config"
	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

var (
	// errBucketPermissionDenied is returned if the bucket ACL denies access to the requested
	// bucket for the associated RPC.
	errBucketPermissionDenied = status.Error(
		codes.PermissionDenied,
		"backend: access to bucket denied by policy",
	)
)

// Permission is a Backend that applies a bucket-based ACL with per-RPC granularity.
type Permission struct {
	permissions []*config.Permission
	Backend
}

// NewPermission creates a new Permission with the specified bucket ACL.
func NewPermission(permissions []*config.Permission, backend Backend) Backend {
	return &Permission{
		permissions: permissions,
		Backend:     backend,
	}
}

// HeadBucket defers to the underlying backend only when permitted by the bucket ACL, and populates
// the permissions ACL in the metadata response.
func (p *Permission) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	var allowed bool
	var permission *config.Permission

	for _, permission = range p.permissions {
		if permission.Pattern.MatchString(request.Resource.Bucket) {
			allowed = permission.RPC.HeadBucket
			break
		}
	}

	if !allowed {
		return nil, errBucketPermissionDenied
	}

	response, err := p.Backend.HeadBucket(ctx, request)
	if err != nil {
		return nil, err
	}

	response.Bucket.Permissions = &common.Permissions{
		HeadBucket:      permission.RPC.HeadBucket,
		HeadObject:      permission.RPC.HeadObject,
		GetObject:       permission.RPC.GetObject,
		StreamGetObject: permission.RPC.StreamGetObject,
		PutObject:       permission.RPC.PutObject,
		StreamPutObject: permission.RPC.StreamPutObject,
		DeleteObject:    permission.RPC.DeleteObject,
		ListBuckets:     permission.RPC.ListBuckets,
		ListObjects:     permission.RPC.ListObjects,
	}

	return response, nil
}

// HeadObject defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	var allowed bool

	for _, permission := range p.permissions {
		if permission.Pattern.MatchString(request.Resource.Bucket) {
			allowed = permission.RPC.HeadObject
			break
		}
	}

	if !allowed {
		return nil, errBucketPermissionDenied
	}

	return p.Backend.HeadObject(ctx, request)
}

// GetObject defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	var allowed bool

	for _, permission := range p.permissions {
		if permission.Pattern.MatchString(request.Resource.Bucket) {
			allowed = permission.RPC.GetObject
			break
		}
	}

	if !allowed {
		return nil, errBucketPermissionDenied
	}

	return p.Backend.GetObject(ctx, request)
}

// StreamGetObject defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	var allowed bool

	errs := make(chan error, 1)

	for _, permission := range p.permissions {
		if permission.Pattern.MatchString(request.Request.Resource.Bucket) {
			allowed = permission.RPC.StreamGetObject
			break
		}
	}

	if !allowed {
		errs <- errBucketPermissionDenied
		return nil, errs
	}

	return p.Backend.StreamGetObject(ctx, request)
}

// PutObject defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	var allowed bool

	for _, permission := range p.permissions {
		if permission.Pattern.MatchString(request.Resource.Bucket) {
			allowed = permission.RPC.PutObject
			break
		}
	}

	if !allowed {
		return nil, errBucketPermissionDenied
	}

	return p.Backend.PutObject(ctx, request)
}

// StreamPutObject defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)
	responses := make(chan *service.PutObjectStreamResponse)
	errs := make(chan error)

	go func() {
		for request := range stream {
			var allowed bool

			for _, permission := range p.permissions {
				if permission.Pattern.MatchString(request.Request.Resource.Bucket) {
					allowed = permission.RPC.StreamPutObject
					break
				}
			}

			if !allowed {
				errs <- errBucketPermissionDenied
				break
			}

			proxy <- request
		}

		close(proxy)
	}()

	go func() {
		response, err := p.Backend.StreamPutObject(ctx, proxy)
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

// DeleteObject defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	var allowed bool

	for _, permission := range p.permissions {
		if permission.Pattern.MatchString(request.Resource.Bucket) {
			allowed = permission.RPC.DeleteObject
			break
		}
	}

	if !allowed {
		return nil, errBucketPermissionDenied
	}

	return p.Backend.DeleteObject(ctx, request)
}

// ListBuckets defers to the underlying backend to obtain a list of buckets and filters out those
// that are disallowed by bucket ACL.
func (p *Permission) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	var buckets []*service.HeadBucketResponse

	response, err := p.Backend.ListBuckets(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, bucket := range response.Buckets {
		// Include in the response only those buckets that are readable based on its
		// ListBuckets RPC permission flag. Buckets encountered without this permission are
		// transparently omitted.

		for _, permission := range p.permissions {
			if permission.Pattern.MatchString(bucket.Bucket.Name) && permission.RPC.ListBuckets {
				bucket.Bucket.Permissions = &common.Permissions{
					HeadBucket:      permission.RPC.HeadBucket,
					HeadObject:      permission.RPC.HeadObject,
					GetObject:       permission.RPC.GetObject,
					StreamGetObject: permission.RPC.StreamGetObject,
					PutObject:       permission.RPC.PutObject,
					StreamPutObject: permission.RPC.StreamPutObject,
					DeleteObject:    permission.RPC.DeleteObject,
					ListBuckets:     permission.RPC.ListBuckets,
					ListObjects:     permission.RPC.ListObjects,
				}

				buckets = append(buckets, bucket)
				break
			}
		}
	}

	return &service.ListBucketsResponse{Buckets: buckets}, nil
}

// ListObjects defers to the underlying backend only when permitted by the bucket ACL.
func (p *Permission) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	var allowed bool

	for _, permission := range p.permissions {
		if permission.Pattern.MatchString(request.Resource.Bucket) {
			allowed = permission.RPC.ListObjects
			break
		}
	}

	if !allowed {
		return nil, errBucketPermissionDenied
	}

	return p.Backend.ListObjects(ctx, request)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (p *Permission) Descriptor() *common.Node {
	return &common.Node{
		Name: "permission",
		Children: map[string]*common.Node_Value{
			"backend": {Child: &common.Node_Value_Node{Node: p.Backend.Descriptor()}},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (p *Permission) String() string {
	return schemas.MarshalDescriptor(p.Descriptor())
}
