package backend

import (
	"context"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"unistore/internal/config"
	"unistore/internal/meta"
	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Metadata is a Backend that automatically injects metadata attributes for inbound put requests.
type Metadata struct {
	server      string
	compression common.Compression
	encryption  common.Encryption
	Backend
}

// NewMetadata creates a new Metadata, with some properties derived automatically from the provided
// backend configuration.
func NewMetadata(compression *config.Compression, encryption *config.Encryption, backend Backend) (Backend, error) {
	var err error

	m := &Metadata{Backend: backend}

	m.server, err = os.Hostname()
	if err != nil {
		return nil, err
	}

	if compression != nil {
		m.compression = common.Compression(common.Compression_value[strings.ToUpper(compression.Algorithm)])
	}

	if encryption != nil {
		m.encryption = common.Encryption(common.Encryption_value[strings.ToUpper(encryption.Mechanism)])
	}

	return m, nil
}

// PutObject derives server-side attributes, attaches it to the inbound request, and passes the
// request to the underlying backend.
func (m *Metadata) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	attributes := &common.Attributes{
		Server:      m.server,
		Version:     meta.Version,
		Timestamp:   timestamppb.New(time.Now()),
		Backend:     request.Resource.Backend,
		Compression: m.compression,
		Encryption:  m.encryption,
		Size:        uint64(len(request.Data)),
		Extra:       make(map[string]string),
	}

	if request.Attributes != nil {
		proto.Merge(attributes, request.Attributes)
	}

	request.Attributes = attributes

	return m.Backend.PutObject(ctx, request)
}

// StreamPutObject derives server-side attributes, attaches it to all inbound requests in the
// stream, and passes the stream to the underlying backend.
func (m *Metadata) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	proxy := make(chan *service.PutObjectStreamRequest)

	go func() {
		var size int

		for request := range stream {
			attributes := &common.Attributes{
				Server:      m.server,
				Version:     meta.Version,
				Timestamp:   timestamppb.New(time.Now()),
				Backend:     request.Request.Resource.Backend,
				Compression: m.compression,
				Encryption:  m.encryption,
				Extra:       make(map[string]string),
			}

			size += len(request.Request.Data)
			attributes.Size = uint64(size)

			if request.Request.Attributes != nil {
				proto.Merge(attributes, request.Request.Attributes)
			}

			request.Request.Attributes = attributes

			proxy <- request
		}

		close(proxy)
	}()

	return m.Backend.StreamPutObject(ctx, proxy)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (m *Metadata) Descriptor() *common.Node {
	return &common.Node{
		Name: "metadata",
		Children: map[string]*common.Node_Value{
			"server":  {Child: &common.Node_Value_Value{Value: m.server}},
			"backend": {Child: &common.Node_Value_Node{Node: m.Backend.Descriptor()}},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (m *Metadata) String() string {
	return schemas.MarshalDescriptor(m.Descriptor())
}
