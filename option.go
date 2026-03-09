package goerr

import "maps"

// Option defines a functional option for modifying the Error struct.
type Option func(*Error)

// WithSpec sets the business specification of the error.
func WithSpec(spec Spec) Option {
	return func(e *Error) {
		e.spec = spec
	}
}

// WithFields adds multiple fields to the current error level.
func WithFields(fields map[string]any) Option {
	return func(e *Error) {
		maps.Copy(e.fields, fields)
	}
}

// WithField adds a field to the current error level.
func WithField(key string, value any) Option {
	return func(e *Error) {
		e.fields[key] = value
	}
}

// WithOp adds the operation name to the current error level.
func WithOp(op string) Option {
	return WithField("op", op)
}

// withCause adds the cause of the error to the current error level.
func withCause(err error) Option {
	return func(e *Error) {
		e.cause = err
	}
}
