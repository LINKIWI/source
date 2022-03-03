package backend

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Validation is a Backend that validates the structure and contents of the inbound request before
// passing it to the underlying backend.
type Validation struct {
	Backend
}

// NewValidation creates a new Validation.
func NewValidation(backend Backend) Backend {
	return &Validation{backend}
}

// HeadBucket validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	if err := v.validateResource(request.Resource); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return v.Backend.HeadBucket(ctx, request)
}

// HeadObject validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	if err := v.validateResource(request.Resource); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "backend: object key not defined")
	}

	return v.Backend.HeadObject(ctx, request)
}

// GetObject validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	if err := v.validateResource(request.Resource); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "backend: object key not defined")
	}

	return v.Backend.GetObject(ctx, request)
}

// StreamGetObject validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	errs := make(chan error, 1)

	if request.Request == nil {
		errs <- status.Error(codes.InvalidArgument, "backend: inner request not defined")
		return nil, errs
	} else if err := v.validateResource(request.Request.Resource); err != nil {
		errs <- status.Error(codes.InvalidArgument, err.Error())
		return nil, errs
	}

	return v.Backend.StreamGetObject(ctx, request)
}

// PutObject validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	if err := v.validateResource(request.Resource); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "backend: object key not defined")
	} else if request.Data == nil || len(request.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "backend: object data not defined")
	}

	return v.Backend.PutObject(ctx, request)
}

// StreamPutObject validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)
	responses := make(chan *service.PutObjectStreamResponse)
	errs := make(chan error)

	go func() {
		for request := range stream {
			if request.Request == nil {
				errs <- status.Error(
					codes.InvalidArgument,
					"backend: inner request not defined",
				)
				break
			} else if err := v.validateResource(request.Request.Resource); err != nil {
				errs <- status.Error(codes.InvalidArgument, err.Error())
				break
			} else if request.Request.Key == "" {
				errs <- status.Error(
					codes.InvalidArgument,
					"backend: object key not defined",
				)
				break
			} else if request.Request.Data == nil || len(request.Request.Data) == 0 {
				errs <- status.Error(
					codes.InvalidArgument,
					"backend: object data not defined",
				)
				break
			}

			proxy <- request
		}

		close(proxy)
	}()

	go func() {
		response, err := v.Backend.StreamPutObject(ctx, proxy)
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

// DeleteObject validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	if err := v.validateResource(request.Resource); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "backend: object key not defined")
	}

	return v.Backend.DeleteObject(ctx, request)
}

// ListBuckets validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	if request.Resource == nil {
		return nil, status.Error(codes.InvalidArgument, "backend: resource not defined")
	} else if request.Resource.Backend == common.Backend_EMPTY {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"backend: resource backend not defined: candidates=%v",
			common.Backend_name,
		)
	}

	return v.Backend.ListBuckets(ctx, request)
}

// ListObjects validates the request and passes well-formed requests to the underlying backend.
func (v *Validation) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	if err := v.validateResource(request.Resource); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return v.Backend.ListObjects(ctx, request)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (v *Validation) Descriptor() *common.Node {
	return &common.Node{
		Name: "validation",
		Children: map[string]*common.Node_Value{
			"backend": {Child: &common.Node_Value_Node{Node: v.Backend.Descriptor()}},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (v *Validation) String() string {
	return schemas.MarshalDescriptor(v.Descriptor())
}

// validateResource returns a non-nil error if the resource is invalid.
func (v *Validation) validateResource(resource *common.Resource) error {
	if resource == nil {
		return fmt.Errorf("backend: resource not defined")
	} else if resource.Backend == common.Backend_EMPTY {
		return fmt.Errorf(
			"backend: resource backend not defined: candidates=%v",
			common.Backend_name,
		)
	} else if resource.Bucket == "" {
		return fmt.Errorf(
			"backend: resource bucket name not defined: backend=%v",
			resource.Backend,
		)
	}

	return nil
}
