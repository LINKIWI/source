package mcrpc

import (
	"context"

	"lib.kevinlin.info/mcrpc/protocol"
)

// contextKey is a type alias for keys in the context object passed to all RPC implementations.
type contextKey int

const (
	// ClientConnContextKey identifies the client net.Conn connection associated with a request.
	ClientConnContextKey contextKey = iota

	// ClientRawRequestContextKey identifies the raw command byte buffer sent by the client from
	// which a structured request was successfully parsed.
	ClientRawRequestContextKey
)

// Handler is an interface for a backend implementing the memcache RPC protocol as a frontend.
type Handler interface {
	// Version implements the version command.
	Version(ctx context.Context, request *protocol.VersionRequest) (*protocol.VersionResponse, error)

	// Shutdown implements the shutdown command.
	Shutdown(ctx context.Context, request *protocol.ShutdownRequest) (*protocol.ShutdownResponse, error)

	// Flush implements the flush_all command.
	Flush(ctx context.Context, request *protocol.FlushRequest) (*protocol.FlushResponse, error)

	// Quit implements the quit command.
	Quit(ctx context.Context, request *protocol.QuitRequest) (*protocol.QuitResponse, error)

	// Stats implements the stats command.
	Stats(ctx context.Context, request *protocol.StatsRequest) (*protocol.StatsResponse, error)

	// Watch implements the watch command.
	Watch(ctx context.Context, request *protocol.WatchRequest) (*protocol.WatchResponse, error)

	// Touch implements the touch command.
	Touch(ctx context.Context, request *protocol.TouchRequest) (*protocol.TouchResponse, error)

	// Delete implements the delete command.
	Delete(ctx context.Context, request *protocol.DeleteRequest) (*protocol.DeleteResponse, error)

	// Incr implements the incr command.
	Incr(ctx context.Context, request *protocol.IncrRequest) (*protocol.IncrResponse, error)

	// Decr implements the decr command.
	Decr(ctx context.Context, request *protocol.DecrRequest) (*protocol.DecrResponse, error)

	// Get implements the get command.
	Get(ctx context.Context, request *protocol.GetRequest) (*protocol.GetResponse, error)

	// Gets implements the gets command.
	Gets(ctx context.Context, request *protocol.GetsRequest) (*protocol.GetsResponse, error)

	// Gat implements the gat command.
	Gat(ctx context.Context, request *protocol.GatRequest) (*protocol.GatResponse, error)

	// Gats implements the gats command.
	Gats(ctx context.Context, request *protocol.GatsRequest) (*protocol.GatsResponse, error)

	// Set implements the set command.
	Set(ctx context.Context, request *protocol.SetRequest) (*protocol.SetResponse, error)

	// Add implements the add command.
	Add(ctx context.Context, request *protocol.AddRequest) (*protocol.AddResponse, error)

	// Replace implements the replace command.
	Replace(ctx context.Context, request *protocol.ReplaceRequest) (*protocol.ReplaceResponse, error)

	// Append implements the append command.
	Append(ctx context.Context, request *protocol.AppendRequest) (*protocol.AppendResponse, error)

	// Prepend implements the prepend command.
	Prepend(ctx context.Context, request *protocol.PrependRequest) (*protocol.PrependResponse, error)

	// Cas implements the cas command.
	Cas(ctx context.Context, request *protocol.CasRequest) (*protocol.CasResponse, error)
}
