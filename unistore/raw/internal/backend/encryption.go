package backend

import (
	"io"
	"strings"

	"filippo.io/age"

	"unistore/internal/schemas"
	"unistore/internal/util"
	"unistore/schemas/common"
)

// Encryption is a Backend that offers efficient and transparent in-flight encryption and
// decryption.
type Encryption struct {
	mechanism common.Encryption
	Backend
}

// NewEncryption creates a new Encrypt with the specified private and public keys for decryption and
// encryption, respectively.
func NewEncryption(mechanism string, privateKey string, publicKey string, backend Backend) Backend {
	var encryptor func(io.Writer) (io.WriteCloser, error)
	var decryptor func(io.Reader) (io.ReadCloser, error)

	parsed := common.Encryption(common.Encryption_value[strings.ToUpper(mechanism)])

	switch parsed {
	case common.Encryption_AGE:
		encryptor = func(raw io.Writer) (io.WriteCloser, error) {
			recipient, err := age.ParseX25519Recipient(publicKey)
			if err != nil {
				return nil, err
			}

			return age.Encrypt(raw, recipient)
		}
		decryptor = func(raw io.Reader) (io.ReadCloser, error) {
			identity, err := age.ParseX25519Identity(privateKey)
			if err != nil {
				return nil, err
			}

			plaintext, err := age.Decrypt(raw, identity)
			if err != nil {
				return nil, err
			}

			return io.NopCloser(plaintext), nil
		}

	case common.Encryption_UNENCRYPTED:
		encryptor = func(raw io.Writer) (io.WriteCloser, error) {
			return util.NopWriteCloser(raw), nil
		}
		decryptor = func(raw io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(raw), nil
		}
	}

	e := &Encryption{
		mechanism: parsed,
		Backend:   newIOProcessor(decryptor, encryptor, 0, 0, backend),
	}

	return newInstrumentation("encryption", e)
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (e *Encryption) Descriptor() *common.Node {
	return &common.Node{
		Name: "encryption",
		Children: map[string]*common.Node_Value{
			"mechanism": {
				Child: &common.Node_Value_Value{
					Value: e.mechanism.String(),
				},
			},
			"backend": {
				Child: &common.Node_Value_Node{
					Node: e.Backend.Descriptor(),
				},
			},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (e *Encryption) String() string {
	return schemas.MarshalDescriptor(e.Descriptor())
}
