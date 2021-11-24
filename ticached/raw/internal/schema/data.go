package schema

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/spaolacci/murmur3"
)

var (
	// ErrShortData is returned when the data passed to UnmarshalBinary is too small to possibly
	// accommodate the full binary format of either a key or value.
	ErrShortData = errors.New("schema: data length cannot accommodate full schema definition")
	// ErrCorruptedData is returned when an unresolvable discrepancy is encountered while
	// unmarshaling binary data into a structured form.
	ErrCorruptedData = errors.New("schema: value fields are not consistent and may be corrupted")
	// ErrIncompatibleVersion is returned when the declared schema version of the value is not
	// equal to the current schema version recognized by this code.
	ErrIncompatibleVersion = errors.New("schema: value schema version is not compatible")
)

// Key is a structured representation of a key.
type Key struct {
	Salt [8]byte
	Key  string
}

// NewKey creates a structured key from a string key.
func NewKey(key string) (*Key, error) {
	var salt [8]byte

	k := &Key{Salt: salt, Key: key}
	h := murmur3.New64()

	if _, err := h.Write([]byte(key)); err != nil {
		return nil, err
	}

	copy(k.Salt[:], h.Sum(nil))

	return k, nil
}

// MarshalBinary marshals the key into a flat binary format.
func (k *Key) MarshalBinary() ([]byte, error) {
	data := make([]byte, len(k.Salt)+len(k.Key))

	copy(data[:8], k.Salt[:])
	copy(data[8:], k.Key)

	return data, nil
}

// UnmarshalBinary unmarshals binary key data into a structured Key.
func (k *Key) UnmarshalBinary(data []byte) error {
	var salt [8]byte

	if len(data) < 9 {
		return ErrShortData
	}

	copy(salt[:], data[:8])
	k.Salt = salt
	k.Key = string(data[8:])

	return nil
}

// Value is a structured representation of a value.
type Value struct {
	Magic     [4]byte
	Version   uint32
	Reserved  [4]byte
	Kind      ValueKind
	Flags     uint16
	TTL       uint32
	CasUnique uint64
	Data      []byte
}

// NewValue creates a structured value from relevant storage request attributes.
func NewValue(kind ValueKind, flags uint16, ttl uint32, casUnique uint64, data []byte) (*Value, error) {
	switch kind {
	case Simple, Arithmetic:
	default:
		return nil, fmt.Errorf("schema: unknown value kind: kind=%v", kind)
	}

	return &Value{
		Magic:     Magic,
		Version:   Version,
		Reserved:  Reserved,
		Kind:      kind,
		Flags:     flags,
		TTL:       ttl,
		CasUnique: casUnique,
		Data:      data,
	}, nil
}

// MarshalBinary marshals the value into a flat binary format.
func (v *Value) MarshalBinary() ([]byte, error) {
	header := make([]byte, 32)

	copy(header, v.Magic[:])
	binary.BigEndian.PutUint32(header[4:8], v.Version)
	copy(header[8:12], v.Reserved[:])
	binary.BigEndian.PutUint16(header[12:14], uint16(v.Kind))
	binary.BigEndian.PutUint16(header[14:16], v.Flags)
	binary.BigEndian.PutUint32(header[16:20], v.TTL)
	binary.BigEndian.PutUint32(header[20:24], uint32(len(v.Data)))
	binary.BigEndian.PutUint64(header[24:32], v.CasUnique)

	return append(header, v.Data...), nil
}

// UnmarshalBinary unmarshals binary value data into a structured Value.
func (v *Value) UnmarshalBinary(data []byte) error {
	var magic [4]byte
	var reserved [4]byte

	if len(data) < 33 {
		return ErrShortData
	}

	copy(magic[:], data[0:4])
	v.Magic = magic
	if v.Magic != Magic {
		return ErrCorruptedData
	}

	v.Version = binary.BigEndian.Uint32(data[4:8])
	if v.Version != Version {
		return ErrIncompatibleVersion
	}

	copy(reserved[:], data[8:12])
	v.Reserved = reserved
	v.Kind = ValueKind(binary.BigEndian.Uint16(data[12:14]))
	v.Flags = binary.BigEndian.Uint16(data[14:16])
	v.TTL = binary.BigEndian.Uint32(data[16:20])
	size := binary.BigEndian.Uint32(data[20:24])
	v.CasUnique = binary.BigEndian.Uint64(data[24:32])
	v.Data = data[32:]

	if uint32(len(v.Data)) != size {
		return ErrCorruptedData
	}

	return nil
}
