package goerr

import (
	"errors"
	"fmt"
)

// New creates a new Error with a cause error and optional configurations.
func New(cause any, opts ...Option) error {
	// if cause is nil, return nil
	if cause == nil {
		return nil
	}

	// initialize error
	err := &Error{
		fields: make(map[string]any),
	}

	// set cause error
	switch v := cause.(type) {
	case string:
		err.cause = errors.New(v)
	case error:
		err.cause = v
	default:
		return fmt.Errorf("goerr: unsupported main error type %T", v)
	}

	// apply options
	for _, opt := range opts {
		opt(err)
	}

	// return error
	return err
}

// Wrap wraps an error with a cause error and optional configurations.
func Wrap(prev error, context any, opts ...Option) error {
	// if prev is nil, return nil
	if prev == nil {
		return nil
	}

	// extract prev error
	prevErr := FromError(prev)

	// create context error
	ctx := FromError(New(context, opts...))
	// set prev error kind
	ctx.kind = prevErr.Kind()

	// add prev error
	ctx.wrapped = prevErr

	// return context error
	return ctx
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
		cause:  err,
		fields: map[string]any{},
	}
}
