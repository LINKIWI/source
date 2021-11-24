package protocol

import (
	"fmt"
	"time"
)

// Command enumerates all recognized protocol commands.
type Command int

// Enumeration of all recognized commands.
const (
	Unknown Command = iota
	Version
	Shutdown
	Flush
	Stats
	Watch
	Touch
	Delete
	Incr
	Decr
	Get
	Gets
	Gat
	Gats
	Set
	Add
	Replace
	Append
	Prepend
	Cas
)

// StoreStatus enumerates known responses for storage commands.
type StoreStatus int

const (
	// Stored indicates the value was stored.
	Stored StoreStatus = iota
	// NotStored indicates the value was not stored.
	NotStored
	// Exists indicates, for CAS commands, that there is a mismatch in the CAS unique.
	Exists
	// NotFound indicates, for CAS commands, that the key does not exist.
	NotFound
)

// Parser is an interface for a protocol parser.
type Parser interface {
	// Parse parses the command into a structured Request.
	Parse(command []byte) (Request, error)
}

// Request is an interface for request types.
type Request interface {
	// IsNoReply indicates whether the request sets the noreply flag to indicate that the server
	// should not send a response after request processing completes.
	IsNoReply() bool

	fmt.Stringer
}

// Response is an interface for response types.
type Response interface {
	fmt.Stringer
}

// Error is a Response indicating that the client sent an unknown command or one that failed to
// parse.
type Error struct {
	Err error
}

// ClientError is a Response indicating that the request is structurally sound (i.e. parsed
// successfully) but is malformed.
type ClientError struct {
	Err error
}

// ServerError is a Response indicating that a server-side error was encountered while trying to
// serve the request.
type ServerError struct {
	Err error
}

// VersionRequest is a Request for the application version.
type VersionRequest struct{}

// VersionResponse is a Response for the application version.
type VersionResponse struct {
	Version string
}

// ShutdownRequest is a Request to shut down the application.
type ShutdownRequest struct {
	Type string
}

// ShutdownResponse is the Response counterpart to ShutdownRequest.
type ShutdownResponse struct{}

// FlushRequest is a Request to flush the cache of all entries.
type FlushRequest struct {
	Delay time.Duration
}

// FlushResponse is the Response counterpart to FlushRequest.
type FlushResponse struct{}

// StatsRequest is a Request for application statistics.
type StatsRequest struct {
	Type string
}

// StatsResponse is a Response with zero or more statistics items.
type StatsResponse struct {
	Stats []StatsItem
}

// StatsItem is a structured representation of a single statistic, represented as a key-value pair.
type StatsItem struct {
	Name  string
	Value string
}

// WatchRequest is a Request to stream logs from the logger backend.
type WatchRequest struct {
	Loggers []string
}

// WatchResponse is the Response with zero or more log entries.
type WatchResponse struct {
	Logs []LogEntry
}

// LogAttribute is a key-value pair for a single attribute in one log entry.
type LogAttribute struct {
	Key   string
	Value string
}

// LogEntry is one log line from a watcher.
type LogEntry struct {
	Attributes []LogAttribute
}

// TouchRequest is a Request to refresh the TTL of an existing item.
type TouchRequest struct {
	Key        string
	Expiration time.Duration
	NoReply    bool
}

// TouchResponse is the Response counterpart to TouchRequest.
type TouchResponse struct {
	Found bool
}

// DeleteRequest is a Request to delete an item.
type DeleteRequest struct {
	Key     string
	NoReply bool
}

// DeleteResponse is the Response counterpart to DeleteRequest.
type DeleteResponse struct {
	Found bool
}

// IncrRequest is a Request to increment an integer value.
type IncrRequest struct {
	Key     string
	Delta   uint64
	NoReply bool
}

// IncrResponse is the Response counterpart to IncrRequest.
type IncrResponse struct {
	Found bool
	Value uint64
}

// DecrRequest is a Request to decrement an integer value.
type DecrRequest struct {
	Key     string
	Delta   uint64
	NoReply bool
}

// DecrResponse is the Response counterpart to DecrRequest.
type DecrResponse struct {
	Found bool
	Value uint64
}

// GetRequest is a Request to fetch the values associated with one or more keys.
type GetRequest struct {
	Keys []string
}

// GetResponse is the Response counterpart to GetRequest.
type GetResponse struct {
	Values []*Retrieval
}

// GetsRequest is a Request to fetch the values and CAS unique values associated with one or more
// keys.
type GetsRequest struct {
	Keys []string
}

// GetsResponse is the Response counterpart to GetsRequest.
type GetsResponse struct {
	Values []*Retrieval
}

// GatRequest is a Request to get and touch one or more keys.
type GatRequest struct {
	Expiration time.Duration
	Keys       []string
}

// GatResponse is the Response counterpart to GatRequest.
type GatResponse struct {
	Values []*Retrieval
}

// GatsRequest is a Request to get and touch one or more keys with their CAS unique values.
type GatsRequest struct {
	Expiration time.Duration
	Keys       []string
}

// GatsResponse is the Response counterpart to GatsRequest.
type GatsResponse struct {
	Values []*Retrieval
}

// Retrieval describes elements of a value for a get-like request.
type Retrieval struct {
	Key       string
	Flags     uint16
	Size      int
	CasUnique uint64
	Data      []byte
}

// SetRequest is a Request to set a key.
type SetRequest struct {
	Payload *Storage
}

// SetResponse is the Response counterpart to SetRequest.
type SetResponse struct {
	Status StoreStatus
}

// AddRequest is a Request to set a key, but only if it does not already exist.
type AddRequest struct {
	Payload *Storage
}

// AddResponse is the Response counterpart to AddRequest.
type AddResponse struct {
	Status StoreStatus
}

// ReplaceRequest is a Request to replace the value of an existing key.
type ReplaceRequest struct {
	Payload *Storage
}

// ReplaceResponse is the Response counterpart to ReplaceRequest.
type ReplaceResponse struct {
	Status StoreStatus
}

// AppendRequest is a Request to append data to the end of the value associated with an existing
// key.
type AppendRequest struct {
	Payload *Storage
}

// AppendResponse is the Response counterpart to AppendRequest.
type AppendResponse struct {
	Status StoreStatus
}

// PrependRequest is a Request to prepend data to the beginning of the value associated with an
// existing key.
type PrependRequest struct {
	Payload *Storage
}

// PrependResponse is the Response counterpart to PrependRequest.
type PrependResponse struct {
	Status StoreStatus
}

// CasRequest is a Request to store a key, with a unique value for check-and-set semantics.
type CasRequest struct {
	Payload *Storage
}

// CasResponse is the Response counterpart to CasRequest.
type CasResponse struct {
	Status StoreStatus
}

// Storage describes elements of a value for a update-like request.
type Storage struct {
	Key        string
	Flags      uint16
	Expiration time.Duration
	Size       int
	CasUnique  uint64
	Data       []byte
	NoReply    bool
}
