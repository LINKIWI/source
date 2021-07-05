package relay

import (
	"encoding/json"

	"piper/internal/config"
	"piper/internal/tail"
)

var (
	// MessageLineSerializer serializes a tail.Message into its literal line contents only.
	MessageLineSerializer = &messageLineSerializer{}
	// MessageJSONSerializer serializes a tail.Message as a JSON object.
	MessageJSONSerializer = &messageJSONSerializer{}
	// MessageSerializerRegistry provides a map of string identifiers to the corresponding
	// MessageSerializer implementation
	MessageSerializerRegistry = map[config.Serializer]MessageSerializer{
		"":                    MessageLineSerializer,
		config.SerializerLine: MessageLineSerializer,
		config.SerializerJSON: MessageJSONSerializer,
	}
)

type messageLineSerializer struct{}

type messageJSONSerializer struct{}

// Serialize extracts the line from the tail event message and provides it as-is.
func (m *messageLineSerializer) Serialize(message tail.Message) ([]byte, error) {
	return []byte(message.Line), nil
}

// Serialize serializes the message as a JSON object.
func (m *messageJSONSerializer) Serialize(message tail.Message) ([]byte, error) {
	return json.Marshal(message)
}
