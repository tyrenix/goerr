package goerr

import (
	"errors"
	"fmt"
)

// New creates a new Error with a main error and optional configurations.
func New(main any, args ...any) error {
	// if main is nil, return nil
	if main == nil {
		return nil
	}

	// initialize error
	err := &Error{
		fields: make(map[string]any),
	}

	// set main error
	switch v := main.(type) {
	case string:
		err.mainErr = errors.New(v)
	case error:
		err.mainErr = v
	default:
		return fmt.Errorf("goerr: unsupported main error type %T", v)
	}

	// apply options
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			err.wrapped = append(err.wrapped, errors.New(v))
		case error:
			err.wrapped = append(err.wrapped, v)
		case Option:
			v(err)
		}
	}

	// return error
	return err
}

// Wrap wraps an error with a main error and optional configurations.
func Wrap(main any, context any, opts ...Option) error {
	// if main is nil, return nil
	if main == nil {
		return nil
	}

	// init error
	err := &Error{}

	// extract error
	switch v := main.(type) {
	case string:
		err.mainErr = New(v)
	case *Error:
		err.mainErr = v
		err.kind = v.kind
	case error:
		err.mainErr = New(v)
	default:
		return fmt.Errorf("goerr: unsupported main error type %T", v)
	}

	// create context error
	contextErr := FromError(New(context))
	// apply options
	for _, opt := range opts {
		opt(contextErr)
	}

	// append context error
	err.wrapped = append(err.wrapped, contextErr)

	// return error
	return err
}

// FromError converts any error to *Error, preserving context if possible.
func FromError(err error) *Error {
	// if err is nil, return nil
	if err == nil {
		return nil
	}

	// if err is already *Error, return it
	if goErr, ok := err.(*Error); ok {
		return goErr
	}

	// otherwise, create a new *Error
	return &Error{
		mainErr: err,
		fields:  make(map[string]any),
	}
}
