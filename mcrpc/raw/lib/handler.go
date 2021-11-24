package lib

import (
	"context"
	"errors"

	"lib.kevinlin.info/mcrpc/protocol"
)

// errNotImplemented is a static protocol.ServerError returned for all RPCs.
var errNotImplemented = &protocol.ServerError{
	Err: errors.New("server: handler RPC not implemented"),
}

// NoopHandler is a Handler that leaves all RPCs unimplemented.
type NoopHandler struct{}

// Version is unimplemented.
func (nh *NoopHandler) Version(ctx context.Context, request *protocol.VersionRequest) (*protocol.VersionResponse, error) {
	return nil, errNotImplemented
}

// Shutdown is unimplemented.
func (nh *NoopHandler) Shutdown(ctx context.Context, request *protocol.ShutdownRequest) (*protocol.ShutdownResponse, error) {
	return nil, errNotImplemented
}

// Flush is unimplemented.
func (nh *NoopHandler) Flush(ctx context.Context, request *protocol.FlushRequest) (*protocol.FlushResponse, error) {
	return nil, errNotImplemented
}

// Quit is unimplemented.
func (nh *NoopHandler) Quit(ctx context.Context, request *protocol.QuitRequest) (*protocol.QuitResponse, error) {
	return nil, errNotImplemented
}

// Stats is unimplemented.
func (nh *NoopHandler) Stats(ctx context.Context, request *protocol.StatsRequest) (*protocol.StatsResponse, error) {
	return nil, errNotImplemented
}

// Watch is unimplemented.
func (nh *NoopHandler) Watch(ctx context.Context, request *protocol.WatchRequest) (*protocol.WatchResponse, error) {
	return nil, errNotImplemented
}

// Touch is unimplemented.
func (nh *NoopHandler) Touch(ctx context.Context, request *protocol.TouchRequest) (*protocol.TouchResponse, error) {
	return nil, errNotImplemented
}

// Delete is unimplemented.
func (nh *NoopHandler) Delete(ctx context.Context, request *protocol.DeleteRequest) (*protocol.DeleteResponse, error) {
	return nil, errNotImplemented
}

// Incr is unimplemented.
func (nh *NoopHandler) Incr(ctx context.Context, request *protocol.IncrRequest) (*protocol.IncrResponse, error) {
	return nil, errNotImplemented
}

// Decr is unimplemented.
func (nh *NoopHandler) Decr(ctx context.Context, request *protocol.DecrRequest) (*protocol.DecrResponse, error) {
	return nil, errNotImplemented
}

// Get is unimplemented.
func (nh *NoopHandler) Get(ctx context.Context, request *protocol.GetRequest) (*protocol.GetResponse, error) {
	return nil, errNotImplemented
}

// Gets is unimplemented.
func (nh *NoopHandler) Gets(ctx context.Context, request *protocol.GetsRequest) (*protocol.GetsResponse, error) {
	return nil, errNotImplemented
}

// Gat is unimplemented.
func (nh *NoopHandler) Gat(ctx context.Context, request *protocol.GatRequest) (*protocol.GatResponse, error) {
	return nil, errNotImplemented
}

// Gats is unimplemented.
func (nh *NoopHandler) Gats(ctx context.Context, request *protocol.GatsRequest) (*protocol.GatsResponse, error) {
	return nil, errNotImplemented
}

// Set is unimplemented.
func (nh *NoopHandler) Set(ctx context.Context, request *protocol.SetRequest) (*protocol.SetResponse, error) {
	return nil, errNotImplemented
}

// Add is unimplemented.
func (nh *NoopHandler) Add(ctx context.Context, request *protocol.AddRequest) (*protocol.AddResponse, error) {
	return nil, errNotImplemented
}

// Replace is unimplemented.
func (nh *NoopHandler) Replace(ctx context.Context, request *protocol.ReplaceRequest) (*protocol.ReplaceResponse, error) {
	return nil, errNotImplemented
}

// Append is unimplemented.
func (nh *NoopHandler) Append(ctx context.Context, request *protocol.AppendRequest) (*protocol.AppendResponse, error) {
	return nil, errNotImplemented
}

// Prepend is unimplemented.
func (nh *NoopHandler) Prepend(ctx context.Context, request *protocol.PrependRequest) (*protocol.PrependResponse, error) {
	return nil, errNotImplemented
}

// Cas is unimplemented.
func (nh *NoopHandler) Cas(ctx context.Context, request *protocol.CasRequest) (*protocol.CasResponse, error) {
	return nil, errNotImplemented
}
