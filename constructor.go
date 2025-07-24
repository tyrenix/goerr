package goerr

import (
	"errors"
	"fmt"
	"maps"
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
	case *Error:
		err.mainErr = v.mainErr
		err.wrapped = append(err.wrapped, v.wrapped...)
		maps.Copy(err.fields, v.fields)
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
