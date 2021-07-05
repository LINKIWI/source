package transport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture/internal/test"
)

func TestNewFramedBackendError(t *testing.T) {
	t.Parallel()

	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	_, err := NewFramed(errorBackendFactory, []byte{})
	assert.Error(t, err)
}

func TestNewFramedSendTailer(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	framed, _ := NewFramed(backendFactory, []byte{1, 2})

	n, err := framed.Send([]byte{1, 2, 3})
	assert.Equal(t, 5, n)
	assert.NoError(t, err)
	assert.Equal(t, 1, backend.NumTransmissions())
	assert.Contains(t, backend.Transmissions(), []byte{1, 2, 3, 1, 2})
}

func TestFramedClose(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	framed, _ := NewFramed(backendFactory, []byte{})

	assert.Equal(t, test.ActiveTransport, backend.State())
	assert.NoError(t, framed.Close())
	assert.Equal(t, test.ClosedTransport, backend.State())
}
