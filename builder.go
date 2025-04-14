package goerr

import (
	"errors"
	"fmt"
)

// New creates a new Error instance by wrapping a base error,
// optional additional errors, strings, or configuration options.
func New(base any, opts ...any) error {
	if base == nil || base == "" {
		return nil
	}

	// create new error
	goerr := &Error{}

	// parse base
	switch v := base.(type) {
	case *Error:
		goerr.err = v
		goerr.httpCode = v.httpCode
	case string:
		goerr.err = errors.New(v)
	case error:
		goerr.err = v
	default:
		return fmt.Errorf("goerr: unsupported option type %T", v)
	}

	// // if base is already a *goerr.Error, copy its httpCode
	// if be, ok := base.(*Error); ok {
	// 	goerr.httpCode = be.httpCode
	// }

	// merge errors and apply options
	joined := goerr.err
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

	// set error
	goerr.err = joined

	// return error
	return goerr
}
