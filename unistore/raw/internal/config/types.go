package config

import (
	"time"
)

// Alias describes a canonical name mapping associated with a regular expression pattern.
type Alias struct {
	Pattern *Regex `yaml:"pattern"`
	Alias   string `yaml:"alias"`
}

// Permission describes an RPC authorization policy associated with a regular expression pattern.
type Permission struct {
	Pattern *Regex `yaml:"pattern"`
	RPC     struct {
		HeadBucket      bool `yaml:"head_bucket"`
		HeadObject      bool `yaml:"head_object"`
		GetObject       bool `yaml:"get_object"`
		StreamGetObject bool `yaml:"stream_get_object"`
		PutObject       bool `yaml:"put_object"`
		StreamPutObject bool `yaml:"stream_put_object"`
		DeleteObject    bool `yaml:"delete_object"`
		ListBuckets     bool `yaml:"list_buckets"`
		ListObjects     bool `yaml:"list_objects"`
	} `yaml:"rpc"`
}

// TLSContext describes parameters for TLS origination to storage backends.
type TLSContext struct {
	VerifyPeer   string `yaml:"verify_peer"`
	Certificates []struct {
		Key         string `yaml:"key"`
		Certificate string `yaml:"certificate"`
		Authority   string `yaml:"authority"`
	} `yaml:"certificates"`
}

// Checksum describes per-backend get/put checksum injection and validation configuration.
type Checksum struct {
	Algorithm string `yaml:"algorithm"`
}

// Compression describes per-backend in-flight compression configuration.
type Compression struct {
	Algorithm string `yaml:"algorithm"`
	Level     int    `yaml:"level"`
}

// Encryption describes per-backend in-flight encryption configuration.
type Encryption struct {
	Name       string `yaml:"name"`
	Mechanism  string `yaml:"mechanism"`
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
}

// Buffer describes per-backend get/put buffering configuration.
type Buffer struct {
	GetChunkSize uint64 `yaml:"get_chunk_size"`
	PutChunkSize uint64 `yaml:"put_chunk_size"`
}

// Threshold describes per-backend stream get/put size threshold configuration.
type Threshold struct {
	MinStreamPutSize int `yaml:"min_stream_put_size"`
	MaxUnaryPutSize  int `yaml:"max_unary_put_size"`
}

// Connection describes connection options for remote backends.
type Connection struct {
	Identity       string        `yaml:"identity"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	TLSContext     *TLSContext   `yaml:"tls_context"`
}

// Composite describes options for the composite object storage backend, which in itself combines
// other declared storage backends.
type Composite struct {
	ReadDispatch  string   `yaml:"read_dispatch"`
	WriteDispatch string   `yaml:"write_dispatch"`
	Backends      []string `yaml:"backends"`
}

// Disk describes options for the disk object storage backend.
type Disk struct {
	Root              string        `yaml:"root"`
	BucketAliases     []*Alias      `yaml:"bucket_aliases"`
	BucketPermissions []*Permission `yaml:"bucket_permissions"`
	Buffer            *Buffer       `yaml:"buffer"`
	Checksum          *Checksum     `yaml:"checksum"`
	Compression       *Compression  `yaml:"compression"`
	Encryption        *Encryption   `yaml:"encryption"`
}

// B2 describes options for the Backblaze B2 object storage backend.
type B2 struct {
	ApplicationKeyID  string        `yaml:"application_key_id"`
	ApplicationKey    string        `yaml:"application_key"`
	Connection        Connection    `yaml:"connection"`
	BucketAliases     []*Alias      `yaml:"bucket_aliases"`
	BucketPermissions []*Permission `yaml:"bucket_permissions"`
	Buffer            *Buffer       `yaml:"buffer"`
	Threshold         *Threshold    `yaml:"threshold"`
	Checksum          *Checksum     `yaml:"checksum"`
	Compression       *Compression  `yaml:"compression"`
	Encryption        *Encryption   `yaml:"encryption"`
}

// Unistore describes options for the Unistore object storage backend.
type Unistore struct {
	Address           string        `yaml:"address"`
	Authority         string        `yaml:"authority"`
	Connection        Connection    `yaml:"connection"`
	Backend           string        `yaml:"backend"`
	BucketAliases     []*Alias      `yaml:"bucket_aliases"`
	BucketPermissions []*Permission `yaml:"bucket_permissions"`
	Buffer            *Buffer       `yaml:"buffer"`
	Threshold         *Threshold    `yaml:"threshold"`
	Checksum          *Checksum     `yaml:"checksum"`
	Compression       *Compression  `yaml:"compression"`
	Encryption        *Encryption   `yaml:"encryption"`
}

// Storage defines available object storage backends.
type Storage struct {
	Composite *Composite `yaml:"composite"`
	Disk      *Disk      `yaml:"disk"`
	Unistore  *Unistore  `yaml:"unistore"`
	B2        *B2        `yaml:"b2"`
}

// RPC describes gRPC-specific server options.
type RPC struct {
	MaxRecvMessageSize   int `yaml:"max_recv_message_size"`
	MaxSendMessageSize   int `yaml:"max_send_message_size"`
	MaxConcurrentStreams int `yaml:"max_concurrent_streams"`
}

// Listener describes a single network listener for the gRPC server.
type Listener struct {
	Address              *Address `yaml:"address"`
	MaxActiveConnections int      `yaml:"max_active_connections"`
}

// Server describes server configuration.
type Server struct {
	Listener Listener `yaml:"listener"`
	RPC      RPC      `yaml:"rpc"`
	Storage  Storage  `yaml:"storage"`
}

// Store describes client connection parameters for a Unistore server instance.
type Store struct {
	Name       string     `yaml:"name"`
	Address    string     `yaml:"address"`
	Authority  string     `yaml:"authority"`
	Connection Connection `yaml:"connection"`
	Backend    string     `yaml:"backend"`
}

// Client describes client configuration.
type Client struct {
	Stores []Store `yaml:"stores"`
	RPC    RPC     `yaml:"rpc"`
}

// Metrics describes metrics client configuration.
type Metrics struct {
	Address    string `yaml:"address"`
	Prefix     string `yaml:"prefix"`
	Proxy      string `yaml:"proxy"`
	Serializer string `yaml:"serializer"`
}

// Log describes application logging configuration.
type Log struct {
	Level   string   `yaml:"level"`
	Outputs []string `yaml:"outputs"`
}

// Errors describes application error tracking and capture configuration.
type Errors struct {
	SentryDSN string `yaml:"sentry_dsn"`
}

// Application describes application-level configuration.
type Application struct {
	Metrics *Metrics `yaml:"metrics"`
	Log     *Log     `yaml:"log"`
	Errors  *Errors  `yaml:"errors"`
}

// Config is the top-level configuration.
type Config struct {
	Application Application `yaml:"application"`
	Server      *Server     `yaml:"server,omitempty"`
	Client      *Client     `yaml:"client,omitempty"`
}
