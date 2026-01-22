package goerr

import (
	"fmt"
)

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
