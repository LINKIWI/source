package server

import (
	"context"
	"time"

	"github.com/tikv/client-go/v2/rawkv"
	"lib.kevinlin.info/mcrpc"
	"lib.kevinlin.info/mcrpc/lib"
	"lib.kevinlin.info/mcrpc/protocol"

	"ticached/internal/meta"
	"ticached/internal/schema"
)

const (
	// kvTimeout is the timeout for all TiKV RPCs.
	kvTimeout = 1 * time.Second
)

// TiKVHandler is a mcrpc.Handler backed by TiKV.
type TiKVHandler struct {
	kv *rawkv.Client
	lib.NoopHandler
}

// NewTiKVHandler creates a new TiKVHandler.
func NewTiKVHandler(kv *rawkv.Client) mcrpc.Handler {
	return &TiKVHandler{kv: kv}
}

// Version reports the current application version.
func (t *TiKVHandler) Version(ctx context.Context, request *protocol.VersionRequest) (*protocol.VersionResponse, error) {
	return &protocol.VersionResponse{Version: meta.Version}, nil
}

// Get performs a non-transactional lookup of the requested key.
func (t *TiKVHandler) Get(ctx context.Context, request *protocol.GetRequest) (*protocol.GetResponse, error) {
	kvctx, cancel := context.WithTimeout(ctx, kvTimeout)
	defer cancel()

	keys := make([][]byte, len(request.Keys))
	for idx, rawKey := range request.Keys {
		key, err := schema.NewKey(rawKey)
		if err != nil {
			return nil, &protocol.ServerError{Err: err}
		}

		keyb, err := key.MarshalBinary()
		if err != nil {
			return nil, &protocol.ServerError{Err: err}
		}

		keys[idx] = keyb
	}

	values, err := t.kv.BatchGet(kvctx, keys)
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	retrievals := make([]*protocol.Retrieval, 0, len(request.Keys))
	for idx, rawValue := range values {
		if len(rawValue) == 0 {
			// Key not found; skip this retrieval data block.
			continue
		}

		value := &schema.Value{}
		if err := value.UnmarshalBinary(rawValue); err != nil {
			return nil, &protocol.ServerError{Err: err}
		}

		retrievals = append(retrievals, &protocol.Retrieval{
			Key:       request.Keys[idx],
			Flags:     value.Flags,
			Size:      len(value.Data),
			CasUnique: value.CasUnique,
			Data:      value.Data,
		})
	}

	return &protocol.GetResponse{Values: retrievals}, nil
}

// Set performs a non-transactional write of a key-value pair, optionally with a finite TTL.
func (t *TiKVHandler) Set(ctx context.Context, request *protocol.SetRequest) (*protocol.SetResponse, error) {
	kvctx, cancel := context.WithTimeout(ctx, kvTimeout)
	defer cancel()

	key, err := schema.NewKey(request.Payload.Key)
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	keyb, err := key.MarshalBinary()
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	value, err := schema.NewValue(
		schema.Simple,
		request.Payload.Flags,
		uint32(request.Payload.Expiration.Seconds()),
		request.Payload.CasUnique,
		request.Payload.Data,
	)
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	valueb, err := value.MarshalBinary()
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	if request.Payload.Expiration > 0 {
		if err := t.kv.PutWithTTL(kvctx, keyb, valueb, uint64(request.Payload.Expiration.Seconds())); err != nil {
			return nil, &protocol.ServerError{Err: err}
		}
	} else {
		if err := t.kv.Put(kvctx, keyb, valueb); err != nil {
			return nil, &protocol.ServerError{Err: err}
		}
	}

	return &protocol.SetResponse{Status: protocol.Stored}, nil
}

// Delete performs a non-transactional delete of the requested key.
func (t *TiKVHandler) Delete(ctx context.Context, request *protocol.DeleteRequest) (*protocol.DeleteResponse, error) {
	kvctx, cancel := context.WithTimeout(ctx, kvTimeout)
	defer cancel()

	key, err := schema.NewKey(request.Key)
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	keyb, err := key.MarshalBinary()
	if err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	if err := t.kv.Delete(kvctx, keyb); err != nil {
		return nil, &protocol.ServerError{Err: err}
	}

	// TiKV client API doesn't expose whether the key previously existed. For purposes of the
	// client API, always assume an existing key was deleted.
	return &protocol.DeleteResponse{Found: true}, nil
}
