package relay

import (
	"piper/internal/tail"
)

// MessageSerializer can serialize tail.Message objects.
type MessageSerializer interface {
	// Serialize serializes a tail.Message into a bytestream.
	Serialize(message tail.Message) ([]byte, error)
}
