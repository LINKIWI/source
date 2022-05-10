package transport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture/internal/test"
)

func TestNewAsyncBackendError(t *testing.T) {
	t.Parallel()

	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	_, err := NewAsync(errorBackendFactory, 1)
	assert.Error(t, err)
}

func TestNewAsyncValidationError(t *testing.T) {
	t.Parallel()

	backendFactory := func() (Transport, error) { return test.NewMockTransport(), nil }

	_, err := NewAsync(backendFactory, 0)
	assert.Error(t, err)

	_, err = NewAsync(backendFactory, -1)
	assert.Error(t, err)

	async, err := NewAsync(backendFactory, 1)
	assert.NotNil(t, async)
	assert.NoError(t, err)
}

func TestAsyncSendFlush(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	async, _ := NewAsync(backendFactory, 2)

	n, err := async.Send([]byte{1})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)

	n, err = async.Send([]byte{2})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)

	err = async.Close()
	assert.NoError(t, err)

	assert.Equal(t, 2, backend.NumTransmissions())
}

func TestAsyncSendFullQueue(t *testing.T) {
	t.Parallel()

	backend := test.NewStalledTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	async, _ := NewAsync(backendFactory, 1)

	// Queued and immediately dequeued by the consumer, but stalled during send
	n, err := async.Send([]byte{1})
	assert.Equal(t, 1, n)
	assert.NoError(t, err)

	// Failure to queue due to queue capacity exhaustion
	n, err = async.Send([]byte{2})
	assert.Equal(t, 0, n)
	assert.Error(t, err)
}

func TestAsyncSendClosing(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	async, _ := NewAsync(backendFactory, 2)

	err := async.Close()
	assert.NoError(t, err)

	n, err := async.Send([]byte{1})
	assert.Equal(t, 0, n)
	assert.Error(t, err)
}

func TestAsyncClose(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	async, _ := NewAsync(backendFactory, 2)

	assert.Equal(t, test.ActiveTransport, backend.State())
	assert.NoError(t, async.Close())
	assert.Equal(t, test.ClosedTransport, backend.State())
}
