package goerr

import "errors"

// New creates a new Error instance by wrapping a base error,
// optional additional errors, strings, or configuration options.
func New(base error, opts ...any) error {
	if base == nil {
		return nil
	}

	goerr := &Error{}

	// If base is already a *goerr.Error, copy its httpCode
	if be, ok := base.(*Error); ok {
		goerr.httpCode = be.httpCode
	}

	// Merge errors and apply options
	joined := base
	for _, opt := range opts {
		switch v := opt.(type) {
		case string:
			joined = errors.Join(joined, errors.New(v))
		case error:
			joined = errors.Join(joined, v)
		case Option:
			v(goerr)
		default:
			return nil
		}
	}

	goerr.err = joined
	return goerr
}
