package goerr

import (
	"fmt"
	"strings"
)

// Wrapped returns the list of wrapped errors.
func (e *Error) Wrapped() []error {
	return e.wrapped
}

// Fields returns the map of custom fields.
func (e *Error) Fields() map[string]any {
	return e.fields
}

// GetField returns the value of a field by key, with a boolean indicating if it exists.
func (e *Error) GetField(key string) (any, bool) {
	if e == nil {
		return nil, false
	}
	v, ok := e.fields[key]
	return v, ok
}

// HTTPCode returns the associated HTTP status code.
func (e *Error) HTTPCode() int {
	if v, ok := e.GetField(fieldHTTPCode); ok {
		if httpCode, ok := v.(int); ok {
			return httpCode
		}
	}
	return 0
}

// Format implements fmt.Formatter for custom formatting.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		var parts []string
		if e.mainErr != nil {
			parts = append(parts, e.mainErr.Error())
		}
		for _, w := range e.wrapped {
			parts = append(parts, w.Error())
		}
		fmt.Fprint(s, strings.Join(parts, ": "))
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	case 's':
		fmt.Fprint(s, e.Error())
	}
}
