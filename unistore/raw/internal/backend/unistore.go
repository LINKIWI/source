package backend

import (
	"context"
	"fmt"
	"io"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"unistore/internal/config"
	"unistore/internal/meta"
	"unistore/internal/schemas"
	"unistore/pkg/client/unistore"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Unistore is a Backend that proxies to another Unistore server. Delightfully meta!
type Unistore struct {
	address    string
	authority  string
	backend    common.Backend
	connection config.Connection
	client     *unistore.Client
}

// NewUnistore creates a new Unistore backend pointed at the specified address and authority,
// overriding the target backend alias for outbound requests to the proxied server.
func NewUnistore(cfg *config.Unistore, rpc config.RPC) (Backend, error) {
	var cancel context.CancelFunc

	instance := "unknown"
	if cfg.Connection.Identity != "" {
		instance = cfg.Connection.Identity
	}
	clientID := fmt.Sprintf(
		"unistore/%s (instance:%s; authority:%s)",
		meta.Version,
		instance,
		cfg.Authority,
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithAuthority(cfg.Authority),
		grpc.WithUserAgent(clientID),
	}

	if rpc.MaxRecvMessageSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(rpc.MaxRecvMessageSize)))
	}
	if rpc.MaxSendMessageSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(rpc.MaxSendMessageSize)))
	}

	ctx := context.Background()
	if cfg.Connection.ConnectTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.Connection.ConnectTimeout)
		defer cancel()
	}

	client, err := unistore.New(ctx, cfg.Address, opts...)
	if err != nil {
		return nil, err
	}

	bid, ok := common.Backend_value[strings.ToUpper(cfg.Backend)]
	if !ok {
		return nil, fmt.Errorf(
			"backend: unknown Unistore proxy backend alias: backend=%s",
			cfg.Backend,
		)
	}

	u := &Unistore{
		address:    cfg.Address,
		authority:  cfg.Authority,
		backend:    common.Backend(bid),
		connection: cfg.Connection,
		client:     client,
	}

	return newInstrumentation("unistore", u), nil
}

// HeadBucket invokes the gRPC HeadBucket RPC.
func (u *Unistore) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.HeadBucketRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.HeadBucket(ctx, proxyRequest)
}

// HeadObject invokes the gRPC HeadObject RPC.
func (u *Unistore) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.HeadObjectRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.HeadObject(ctx, proxyRequest)
}

// GetObject invokes the gRPC GetObject RPC.
func (u *Unistore) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.GetObjectRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.GetObject(ctx, proxyRequest)
}

// StreamGetObject invokes the gRPC StreamGetObject RPC and adapts the response stream to a channel.
func (u *Unistore) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	responses := make(chan *service.GetObjectStreamResponse)
	errs := make(chan error, 1)

	proxyRequest := proto.Clone(request).(*service.GetObjectStreamRequest)
	proxyRequest.Request.Resource.Backend = u.backend

	stream, err := u.client.StreamGetObject(ctx, proxyRequest)
	if stream == nil {
		errs <- err
		return nil, errs
	}

	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errs <- err
				}
				break
			}

			responses <- response
		}

		close(responses)
	}()

	return responses, errs
}

// PutObject invokes the gRPC PutObject RPC.
func (u *Unistore) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.PutObjectRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.PutObject(ctx, proxyRequest)
}

// StreamPutObject invokes the gRPC StreamPutObject RPC and adapts the request stream to a channel.
func (u *Unistore) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	c, err := u.client.StreamPutObject(ctx)
	if err != nil {
		return nil, err
	}

	for request := range stream {
		proxyRequest := proto.Clone(request).(*service.PutObjectStreamRequest)
		proxyRequest.Request.Resource.Backend = u.backend

		if err := c.Send(proxyRequest); err != nil {
			break
		}
	}

	return c.CloseAndRecv()
}

// DeleteObject invokes the gRPC DeleteObject RPC.
func (u *Unistore) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.DeleteObjectRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.DeleteObject(ctx, proxyRequest)
}

// ListBuckets invokes the gRPC ListBuckets RPC.
func (u *Unistore) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.ListBucketsRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.ListBuckets(ctx, proxyRequest)
}

// ListObjects invokes the gRPC ListObjects RPC.
func (u *Unistore) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	var cancel context.CancelFunc

	if u.connection.RequestTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, u.connection.RequestTimeout)
		defer cancel()
	}

	proxyRequest := proto.Clone(request).(*service.ListObjectsRequest)
	proxyRequest.Resource.Backend = u.backend

	return u.client.ListObjects(ctx, proxyRequest)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (u *Unistore) Descriptor() *common.Node {
	return &common.Node{
		Name: "unistore",
		Children: map[string]*common.Node_Value{
			"address":   {Child: &common.Node_Value_Value{Value: u.address}},
			"authority": {Child: &common.Node_Value_Value{Value: u.authority}},
			"backend":   {Child: &common.Node_Value_Value{Value: u.backend.String()}},
		},
	}
}

// Close closes the underlying Unistore gRPC client.
func (u *Unistore) Close() error {
	return u.client.Close()
}

// String returns a human-consumable representation of this backend.
func (u *Unistore) String() string {
	return schemas.MarshalDescriptor(u.Descriptor())
}
