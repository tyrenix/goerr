package goerr

import "errors"

// Error represents a custom error with a main error, wrapped errors, and fields.
type Error struct {
	cause error
	// kind represents the business classification of the error.
	// It is set only on creation (New) and is inherited across wrapping.
	// Wrap MUST NOT override kind.
	kind    error
	wrapped error
	fields  map[string]any
}

// Error returns the main error message.
func (e *Error) Error() string {
	// if error is nil, return empty string
	if e == nil {
		return ""
	}

	// get stack
	st := stack(e)
	if len(st) == 0 {
		return ""
	}

	// return last error
	if last, ok := st[len(st)-1].(*Error); ok {
		// if not exists cause return empty string
		if last.cause == nil {
			return ""
		}
		// return cause
		return last.cause.Error()
	}

	// return last error
	return st[len(st)-1].Error()
}

// Unwrap returns all wrapped errors for compatibility with Go 1.20+ errors.Is/As.
func (e *Error) Unwrap() error {
	// if error is nil, return nil
	if e == nil {
		return nil
	}

	// return wrapped error
	return e.wrapped
}

// Is returns true if the target error is the same as the cause or wrapped error.
func (e *Error) Is(target error) bool {
	// if error is nil, return false
	if e == nil || target == nil {
		return false
	}

	// if target equal cause, return true
	if e.cause != nil && errors.Is(e.cause, target) {
		return true
	}

	// if target equal , return true
	if errors.Is(e.cause, target) {
		return true
	}

	return false
}
