package supercharged

import (
	"fmt"
	"net/http"
)

// Error includes additional metadata for a Supercharged error response.
type Error struct {
	Status  int
	Code    int
	Message string
	Data    interface{}
}

// Error returns a string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// WrapError wraps an error with default fields to conform to an Error struct.
func WrapError(err error) *Error {
	return &Error{
		Status:  http.StatusBadRequest,
		Code:    CodeClientUndefined,
		Message: err.Error(),
	}
}
