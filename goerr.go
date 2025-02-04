package goerr

import (
	"errors"
	"fmt"
	"strings"
)

// Error represents a set of errors that have been joined together using
// errors.Join.
type Error struct {
	// err is the underlying error that was created by errors.Join.
	err error
}

// New creates a new Error instance by joining the provided errors.
// The first argument 'err' is a required error, and if it is nil, the function returns nil.
// Additional errors can be provided in 'errs', which are joined to 'err' using errors.Join.
// Each value in 'errs' is converted to an error using extractError before joining.
func New(err error, errs ...any) error {
	// is error nil
	if err == nil {
		return nil
	}

	// add errors to error
	for _, e := range errs {
		err = errors.Join(err, extractError(e))
	}

	// return error
	return &Error{err: err}
}

// Error returns the first error string from the underlying error chain.
// If the underlying error is nil or if there are no errors in the chain,
// it returns an empty string.
func (e *Error) Error() string {
	// is error nil
	if e.err == nil {
		return ""
	}

	// get errors text
	errs := e.unwrap()
	if len(errs) == 0 {
		return ""
	}

	// return error
	return errs[0]
}

// Unwrap returns the underlying error from the Error instance.
// This allows Error to be used with errors.Unwrap and other error handling functions
// that work with wrapped errors.
func (e *Error) Unwrap() error {
	return e.err
}

// unwrap unwraps the underlying error and returns a slice of strings, each of which is an error string.
// If the underlying error is nil, it returns nil.
// If the underlying error is a wrapper, it unwraps the error and returns the string representation of the unwrapped error.
// If the underlying error is not a wrapper, it returns a slice with a single element, which is the string representation of the underlying error.
func (e *Error) unwrap() []string {
	// is error nil
	if e.err == nil {
		return nil
	}

	// unwrap errors and return
	return strings.Split(e.err.Error(), "\n")
}

// Format implements the fmt.Formatter interface.
// It formats the underlying error according to the format specifier 'c'.
// If 'c' is 'v', it formats the error as a string with all errors joined by a colon and a space.
// If 'c' is 'q', it formats the error as a string surrounded by quotes.
// If 'c' is 's', it formats the error as a string without any additional formatting.
func (e *Error) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		fmt.Fprintf(s, "%s", strings.Join(e.unwrap(), ": "))
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	case 's':
		fmt.Fprint(s, e.Error())
	}
}

// extractError extracts an error from a given value. If the value is a string, it is converted to an error using errors.New.
// If the value is an error, it is returned as is. Otherwise, nil is returned.
func extractError(err any) error {
	switch e := err.(type) {
	case string:
		return errors.New(e)
	case error:
		return e
	default:
		return nil
	}
}
