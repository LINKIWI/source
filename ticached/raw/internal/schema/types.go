package schema

var (
	// Version is the current schema version. This version may be incremented as breaking
	// (non-backwards compatible) schema format changes are made.
	Version = uint32(1)

	// Magic is a preamble shared by all value blobs.
	Magic = [4]byte{0x71, 0x60, 0xCA, 0x8E}

	// Reserved is a sequence of bytes reserved for potential future use.
	Reserved = [4]byte{0, 0, 0, 0}
)

// ValueKind is a type alias for various types of values.
type ValueKind uint16

const (
	// Unknown is the default "empty" kind and should be treated as an error condition if
	// encountered in a value payload.
	Unknown ValueKind = iota
	// Simple describes regular value blobs.
	Simple
	// Arithmetic describes numerical value blobs used for arithmetic commands (e.g. increment
	// and decrement). The data payload is a big endian-encoded uint64.
	Arithmetic
)
