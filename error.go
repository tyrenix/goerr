package goerr

import (
	"errors"
)

// Error represents a custom error with a main error, wrapped errors, and fields.
type Error struct {
	mainErr error          // Main error returned by .Error()
	wrapped []error        // Wrapped errors for additional context
	fields  map[string]any // Custom fields (e.g., min_amount: 100)
}

// Error returns the main error message.
func (e *Error) Error() string {
	if e.mainErr == nil {
		return ""
	}
	return e.mainErr.Error()
}

// Unwrap returns the main error for compatibility with errors.Unwrap.
func (e *Error) Unwrap() error {
	return e.mainErr
}

// Is implements errors.Is to check if the error matches a target.
func (e *Error) Is(target error) bool {
	if errors.Is(e.mainErr, target) {
		return true
	}
	for _, w := range e.wrapped {
		if errors.Is(w, target) {
			return true
		}
	}
	return false
}

// As implements errors.As to check if the error can be assigned to a target.
func (e *Error) As(target any) bool {
	if errors.As(e.mainErr, target) {
		return true
	}
	for _, w := range e.wrapped {
		if errors.As(w, target) {
			return true
		}
	}
	return false
}
