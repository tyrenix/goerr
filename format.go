package goerr

import (
	"fmt"
	"sort"
	"strings"
)

// Details returns a detailed string with the primary error, wrapped errors, and fields.
func (e *Error) Details() string {
	var parts []string

	// recursively get details from all unwrapped children
	for _, child := range e.Unwrap() {
		if childGoErr, ok := child.(*Error); ok {
			parts = append(parts, childGoErr.Details())
		} else {
			parts = append(parts, child.Error())
		}
	}

	// add fields from the current error level
	if len(e.fields) > 0 {
		var fieldParts []string
		keys := make([]string, 0, len(e.fields))
		for k := range e.fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, e.fields[k]))
		}
		parts = append(parts, "("+strings.Join(fieldParts, ", ")+")")
	}

	// join all parts
	return strings.Join(parts, ": ")
}

// Format implements fmt.Formatter for custom formatting.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprint(s, e.Details())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	case 's':
		fmt.Fprint(s, e.Error())
	}
}
