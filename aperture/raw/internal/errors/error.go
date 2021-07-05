package errors

import (
	"fmt"
	"strings"
)

// Error provides structured abstractions over errors and supports stacking inner errors.
type Error struct {
	Namespace string
	Message   string
	Tags      []Tag
	Stacked   error
}

// Tag is an optional key-value tag used to provide additional context for an error.
type Tag struct {
	Key   string
	Value interface{}
}

// Error serializes the error as a human-consumable string.
func (e *Error) Error() string {
	var tags []string

	err := fmt.Sprintf("%s: %s", e.Namespace, e.Message)

	for _, tag := range e.Tags {
		tags = append(tags, fmt.Sprintf("%v=%v", tag.Key, tag.Value))
	}

	if len(tags) > 0 {
		err += strings.Join(append([]string{":"}, tags...), " ")
	}

	if e.Stacked != nil {
		err += fmt.Sprintf("\n%s", e.Stacked)
	}

	return err
}

// New creates an error instance.
func New(namespace string, message string, tags ...Tag) *Error {
	return Stack(namespace, message, nil, tags...)
}

// Stack creates a new error instance, stacking an existing inner error.
func Stack(namespace string, message string, err error, tags ...Tag) *Error {
	return &Error{
		Namespace: namespace,
		Message:   message,
		Tags:      tags,
		Stacked:   err,
	}
}
