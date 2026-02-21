package goerr

import (
	"fmt"
	"sort"
	"strings"
)

// Details returns the error details as a string.
func (e *Error) Details() string {
	// init parts
	var parts []string

	// get stack
	stack := stack(e)

	// get last
	last := stack[len(stack)-1]
	// add last to start of stack
	stack = append([]error{last}, stack[:len(stack)-1]...)

	// loop over stack
	for i, err := range stack {
		if ge, ok := err.(*Error); ok {
			parts = append(parts, formatOne(ge, i == 0))
		} else {
			parts = append(parts, err.Error())
		}
	}

	// return formatted parts
	return strings.Join(parts, ": ")
}

// formatOne formats a single error.
func formatOne(e *Error, showKind bool) string {
	// init message
	msg := ""
	if e.cause != nil {
		msg = e.cause.Error()
	}

	// if no fields, return message
	if len(e.fields) == 0 && e.kind == "" {
		return msg
	}

	// init field parts
	var fieldParts []string

	// if kind is set, add it
	if e.kind != "" && showKind {
		fieldParts = append(fieldParts, fmt.Sprintf("kind=%v", e.kind))
	}

	// sort fields
	keys := make([]string, 0, len(e.fields))
	for k := range e.fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// add fields
	for _, k := range keys {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, e.fields[k]))
	}

	// add fields
	if len(fieldParts) > 0 {
		msg += " (" + strings.Join(fieldParts, ", ") + ")"
	}

	// return formatted message
	return msg
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

// stack returns the error stack as a slice of errors.
func stack(err error) []error {
	// init stack
	var stack []error
	for err != nil {
		// add current error
		stack = append(stack, err)
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}

		// get next error
		err = u.Unwrap()
	}

	// return stack
	return stack
}
