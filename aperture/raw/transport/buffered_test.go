package transport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture/internal/test"
)

func TestNewBufferedBackendError(t *testing.T) {
	t.Parallel()

	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	_, err := NewBuffered(errorBackendFactory, 1)
	assert.Error(t, err)
}

func TestNewBufferedValidationError(t *testing.T) {
	t.Parallel()

	backendFactory := func() (Transport, error) { return test.NewMockTransport(), nil }

	_, err := NewBuffered(backendFactory, 0)
	assert.Error(t, err)

	_, err = NewBuffered(backendFactory, -1)
	assert.Error(t, err)

	buffered, err := NewBuffered(backendFactory, 1)
	assert.NotNil(t, buffered)
	assert.NoError(t, err)
}

func TestBufferedBatchedSend(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	buffered, _ := NewBuffered(backendFactory, 2)

	n, err := buffered.Send([]byte{1})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)
	assert.Equal(t, 0, backend.NumTransmissions())

	n, err = buffered.Send([]byte{2, 3, 4})
	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.Equal(t, 2, backend.NumTransmissions())
}

func TestBufferedFlush(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	buffered, _ := NewBuffered(backendFactory, 2)

	n, err := buffered.Send([]byte{1})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)
	assert.Equal(t, 0, backend.NumTransmissions())

	buffered.Close()

	// On close, all queued items are flushed despite not reaching the threshold batch size
	assert.Equal(t, 1, backend.NumTransmissions())
}

func TestBufferedClose(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	buffered, _ := NewBuffered(backendFactory, 2)

	assert.Equal(t, test.ActiveTransport, backend.State())
	assert.NoError(t, buffered.Close())
	assert.Equal(t, test.ClosedTransport, backend.State())
}
