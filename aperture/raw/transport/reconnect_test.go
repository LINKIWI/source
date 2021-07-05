package transport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReconnectingBackendError(t *testing.T) {
	t.Parallel()

	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	_, err := NewReconnecting(errorBackendFactory, 0, 0)
	assert.Error(t, err)
}
