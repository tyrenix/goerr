package goerr

import (
	"fmt"
	"maps"
)

// Wrap wraps an error with a main error and optional configurations.
func Wrap(main any, context any, opts ...Option) error {
	// if main is nil, return nil
	if main == nil {
		return nil
	}

	// initialize source error
	var sourceErr *Error
	// extract error
	switch v := main.(type) {
	case string:
		sourceErr = FromError(New(v))
	case *Error:
		sourceErr = v
	case error:
		sourceErr = FromError(New(v))
	default:
		return fmt.Errorf("goerr: unsupported main error type %T", v)
	}

	// create error
	err := &Error{
		mainErr: sourceErr.mainErr,
		kind:    sourceErr.kind,
		wrapped: make([]error, len(sourceErr.wrapped)),
		fields:  make(map[string]any),
	}

	// copy fields
	maps.Copy(err.fields, sourceErr.fields)
	// copy wrapped
	copy(err.wrapped, sourceErr.wrapped)

	// create context error
	contextErr := FromError(New(context))
	// apply options
	for _, opt := range opts {
		opt(contextErr)
	}

	// prepend context error
	err.wrapped = append([]error{contextErr}, err.wrapped...)

	// return error
	return err
}
