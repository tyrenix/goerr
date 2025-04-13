package goerr

import (
	"fmt"
	"strings"
)

// Error represents a custom error type that can wrap multiple errors
// and carry additional metadata like an HTTP status code.
type Error struct {
	err      error // The underlying error (can be joined errors).
	httpCode int   // Associated HTTP status code.
}

// Error returns the first error message from the joined error list.
func (e *Error) Error() string {
	if e.err == nil {
		return ""
	}
	errs := e.unwrap()
	if len(errs) == 0 {
		return ""
	}
	return errs[0]
}

// Unwrap returns the underlying error, enabling compatibility with errors.Unwrap.
func (e *Error) Unwrap() error {
	return e.err
}

// Format implements fmt.Formatter to allow custom formatting using fmt verbs.
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

// unwrap splits the joined error message into a slice of individual error strings.
func (e *Error) unwrap() []string {
	if e.err == nil {
		return nil
	}
	return strings.Split(e.err.Error(), "\n")
}

// HTTPCode returns the associated HTTP status code.
func (e *Error) HTTPCode() int {
	return e.httpCode
}
