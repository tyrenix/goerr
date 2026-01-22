package goerr

// Error represents a custom error with a main error, wrapped errors, and fields.
type Error struct {
	mainErr error
	kind    error
	wrapped []error
	fields  map[string]any
}

// Error returns the main error message.
func (e *Error) Error() string {
	if e.mainErr == nil {
		return ""
	}
	return e.mainErr.Error()
}

// Unwrap returns all wrapped errors for compatibility with Go 1.20+ errors.Is/As.
func (e *Error) Unwrap() []error {
	// It return a new slice containing mainErr and all from the wrapped slice.
	// This makes the entire chain visible to errors.Is and errors.As.
	all := make([]error, 0, 1+len(e.wrapped))
	if e.mainErr != nil {
		all = append(all, e.mainErr)
	}
	all = append(all, e.wrapped...)
	return all
}
