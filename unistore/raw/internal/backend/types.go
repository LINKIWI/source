package backend

import (
	"context"
	"fmt"

	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Backend describes an object storage implementation.
type Backend interface {
	fmt.Stringer

	// HeadBucket retrieves metadata about a bucket.
	HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error)

	// HeadObject retrieves metadata about an object.
	HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error)

	// GetObject reads an object's full contents (buffered in memory).
	GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error)

	// StreamGetObject is a streaming implementation of GetObject.
	StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error)

	// PutObject writes an object's full contents (buffered in memory).
	PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error)

	// StreamPutObject is a streaming implementation of PutObject.
	StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error)

	// DeleteObject deletes an object.
	DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error)

	// ListBuckets lists all buckets.
	ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error)

	// ListObjects lists objects under a particular prefix.
	ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error)

	// Descriptor provides a structured descriptor of the backend's internal composition tree.
	Descriptor() *common.Node
}
