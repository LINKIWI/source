package backend

import (
	"fmt"

	"unistore/internal/schemas"
	"unistore/schemas/common"
)

// Buffer is a Backend that internally buffers streaming object gets and puts to a desired chunk
// size to the wrapped Backend, independent of the payload sizes sent or received by the client.
// This is particularly useful for backends that have requirements on minimum and maximum chunk
// sizes for streaming I/O.
type Buffer struct {
	getChunkSize uint64
	putChunkSize uint64
	Backend
}

// NewBuffer creates a Buffer with desired get and put chunk sizes. Chunked reads from the backend
// are normalized to at most the get chunk size; chunked writes to the backend are normalized to at
// most the put chunk size.
func NewBuffer(getChunkSize uint64, putChunkSize uint64, backend Backend) Backend {
	return &Buffer{
		getChunkSize: getChunkSize,
		putChunkSize: putChunkSize,
		Backend:      newIOProcessor(nil, nil, getChunkSize, putChunkSize, backend),
	}
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (b *Buffer) Descriptor() *common.Node {
	return &common.Node{
		Name: "buffer",
		Children: map[string]*common.Node_Value{
			"get_chunk_size": {
				Child: &common.Node_Value_Value{
					Value: fmt.Sprintf("%d", b.getChunkSize),
				},
			},
			"put_chunk_size": {
				Child: &common.Node_Value_Value{
					Value: fmt.Sprintf("%d", b.putChunkSize),
				},
			},
			"backend": {
				Child: &common.Node_Value_Node{
					Node: b.Backend.Descriptor(),
				},
			},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (b *Buffer) String() string {
	return schemas.MarshalDescriptor(b.Descriptor())
}
