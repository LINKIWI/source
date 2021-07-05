package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVanillaUntagged(t *testing.T) {
	t.Parallel()

	err := New("namespace", "message")
	assert.Equal(t, "namespace: message", err.Error())
}

func TestVanillaTagged(t *testing.T) {
	t.Parallel()

	err := New("namespace", "message", Tag{"num", 5}, Tag{"bytes", []byte{1, 2, 3}})
	assert.Equal(t, "namespace: message: num=5 bytes=[1 2 3]", err.Error())
}

func TestStacked(t *testing.T) {
	t.Parallel()

	var err *Error

	err = Stack("namespace", "message", fmt.Errorf("stacked"))
	assert.Equal(t, "namespace: message\nstacked", err.Error())

	err = Stack("namespace", "message", New("stacked", "message", Tag{"num", 5}))
	assert.Equal(t, "namespace: message\nstacked: message: num=5", err.Error())

	err = Stack("namespace", "message", New("stacked", "message", Tag{"num", 5}), Tag{"num", 6})
	assert.Equal(t, "namespace: message: num=6\nstacked: message: num=5", err.Error())
}
