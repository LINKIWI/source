package util

import (
	"fmt"
	"strings"
)

// Error provides rich metadata around a standard error.
type Error struct {
	Namespace    string
	Message      string
	Tags         map[string]interface{}
	StackedError error
}

// Unwrap returns the internally stacked error.
func (e *Error) Unwrap() error {
	return e.StackedError
}

// Error provides a string representation of the error.
func (e *Error) Error() string {
	var tags []string

	for key, value := range e.Tags {
		tags = append(tags, fmt.Sprintf("%s=%+v", key, value))
	}

	err := fmt.Sprintf("%s: %s", e.Namespace, e.Message)

	if len(tags) > 0 {
		err += fmt.Sprintf(": %s", strings.Join(tags, " "))
	}

	if e.StackedError != nil {
		err += fmt.Sprintf("\n[stack] %v", e.StackedError)
	}

	return err
}
