package transport

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"lib.kevinlin.info/aperture/internal/test"
)

func TestNewCompositePartialBackendError(t *testing.T) {
	t.Parallel()

	backendFactory := func() (Transport, error) { return test.NewMockTransport(), nil }
	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	_, err := NewComposite([]Factory{backendFactory, errorBackendFactory}, All)

	assert.Error(t, err)
}

func TestNewCompositeFullBackendError(t *testing.T) {
	t.Parallel()

	errorBackendFactory := func() (Transport, error) { return nil, fmt.Errorf("") }

	_, err := NewComposite([]Factory{errorBackendFactory, errorBackendFactory}, All)

	assert.Error(t, err)
}

func TestCompositeEmptyBackends(t *testing.T) {
	t.Parallel()

	composite, err := NewComposite([]Factory{}, All)
	assert.NoError(t, err)

	n, err := composite.Send([]byte{})
	assert.Equal(t, 0, n)
	assert.NoError(t, err)

	err = composite.Close()
	assert.NoError(t, err)
}

func TestCompositeSendFirstConsistencySuccess(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewMockTransport(),
		test.NewErrorTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, First)
	_, err := composite.Send([]byte{})

	assert.NoError(t, err)
}

func TestCompositeSendFirstConsistencyError(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewErrorTransport(),
		test.NewMockTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, First)
	_, err := composite.Send([]byte{})

	assert.Error(t, err)
}

func TestCompositeSendFirstAnyConsistencySuccess(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewErrorTransport(),
		test.NewMockTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, Any)
	_, err := composite.Send([]byte{})

	assert.NoError(t, err)
}

func TestCompositeSendFirstAnyConsistencyError(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewErrorTransport(),
		test.NewErrorTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, Any)
	_, err := composite.Send([]byte{})

	assert.Error(t, err)
}

func TestCompositeSendFirstMajorityConsistencySuccess(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewErrorTransport(),
		test.NewMockTransport(),
		test.NewMockTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, Majority)
	_, err := composite.Send([]byte{})

	assert.NoError(t, err)
}

func TestCompositeSendFirstMajorityConsistencyError(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewErrorTransport(),
		test.NewErrorTransport(),
		test.NewMockTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, Majority)
	_, err := composite.Send([]byte{})

	assert.Error(t, err)
}

func TestCompositeSendFirstAllConsistencySuccess(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewMockTransport(),
		test.NewMockTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, All)
	_, err := composite.Send([]byte{})

	assert.NoError(t, err)
}

func TestCompositeSendFirstAllConsistencyError(t *testing.T) {
	t.Parallel()

	backends := []Transport{
		test.NewMockTransport(),
		test.NewErrorTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, All)
	_, err := composite.Send([]byte{})

	assert.Error(t, err)
}

func TestCompositeClose(t *testing.T) {
	t.Parallel()

	backends := []*test.MockTransport{
		test.NewMockTransport(),
		test.NewMockTransport(),
	}
	factories := []Factory{
		func() (Transport, error) { return backends[0], nil },
		func() (Transport, error) { return backends[1], nil },
	}

	composite, _ := NewComposite(factories, All)

	assert.Equal(t, test.ActiveTransport, backends[0].State())
	assert.Equal(t, test.ActiveTransport, backends[1].State())
	assert.NoError(t, composite.Close())
	assert.Equal(t, test.ClosedTransport, backends[0].State())
	assert.Equal(t, test.ClosedTransport, backends[1].State())
}
