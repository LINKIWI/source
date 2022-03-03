package unistore

import (
	"context"

	"google.golang.org/grpc"

	"unistore/schemas/service"
)

// Client is a Unistore gRPC client wrapper.
type Client struct {
	conn *grpc.ClientConn
	meta service.MetaClient
	service.UnistoreClient
}

// New creates a new Unistore client, which includes dialing the gRPC server.
func New(ctx context.Context, address string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:           conn,
		meta:           service.NewMetaClient(conn),
		UnistoreClient: service.NewUnistoreClient(conn),
	}, nil
}

// Meta provides access to the Meta service client stub.
func (c *Client) Meta() service.MetaClient {
	return c.meta
}

// Close closes the underlying client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
