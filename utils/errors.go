package utils

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is a generic error for things that don't exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidAuth is for authentication events that do not succeed.
	ErrInvalidAuth = errors.New("invalid authentication")

	// ErrRunCanceled is thrown to github when the run is canceled.
	ErrRunCanceled = errors.New("run canceled by user intervention")
)

// WrapError wraps an error with fmt.Error.
func WrapError(err error, message string, args ...interface{}) error {
	return fmt.Errorf("%v: %w", fmt.Sprintf(message, args...), err)
}
