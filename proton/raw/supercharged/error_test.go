package supercharged

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	t.Parallel()

	err := WrapError(errors.New("example"))

	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.Equal(t, CodeClientUndefined, err.Code)
	assert.Equal(t, "example", err.Message)
	assert.Nil(t, err.Data)
}

func TestErrorString(t *testing.T) {
	t.Parallel()

	err := Error{
		Status:  http.StatusInternalServerError,
		Code:    CodeServerUndefined,
		Message: "message",
	}

	assert.Equal(t, "message (-1)", err.Error())
}
