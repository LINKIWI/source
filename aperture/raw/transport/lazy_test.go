package transport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture/internal/test"
)

func TestNewLazy(t *testing.T) {
	t.Parallel()

	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	lazy, err := NewLazy(errorBackendFactory)
	assert.NotNil(t, lazy)
	assert.NoError(t, err)
}

func TestLazySendInitializeError(t *testing.T) {
	t.Parallel()

	initAttempts := 0
	errorBackendFactory := func() (Transport, error) {
		initAttempts++
		return nil, fmt.Errorf("")
	}

	lazy, _ := NewLazy(errorBackendFactory)

	n, err := lazy.Send([]byte{})
	assert.Equal(t, 0, n)
	assert.Error(t, err)
	assert.Equal(t, 1, initAttempts)

	n, err = lazy.Send([]byte{})
	assert.Equal(t, 0, n)
	assert.Error(t, err)
	assert.Equal(t, 2, initAttempts)
}

func TestLazySendInitializeSuccess(t *testing.T) {
	t.Parallel()

	backend := test.NewMockTransport()
	backendFactory := func() (Transport, error) { return backend, nil }
	lazy, _ := NewLazy(backendFactory)

	n, err := lazy.Send([]byte{1, 2, 3})
	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.Equal(t, 1, backend.NumTransmissions())

	n, err = lazy.Send([]byte{1, 2, 3})
	assert.Equal(t, 3, n)
	assert.NoError(t, err)
	assert.Equal(t, 2, backend.NumTransmissions())
}
