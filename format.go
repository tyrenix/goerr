package goerr

import (
	"fmt"
	"sort"
	"strings"
)

// Details returns the error details as a string.
func (e *Error) Details() string {
	return DetailsOf(e)
}

// DetailsOf returns a detailed string for any error.
func DetailsOf(err error) string {
	if err == nil {
		return ""
	}

	var parts []string
	if code, ok := CodeOf(err); ok {
		head := string(code)
		if kind, ok := KindOf(err); ok {
			head += fmt.Sprintf(" (kind=%s)", kind)
		}
		parts = append(parts, head)
	}

	for current := err; current != nil; current = unwrap(current) {
		if goErr, ok := current.(*Error); ok {
			part := formatLevel(goErr)
			if part != "" {
				parts = append(parts, part)
			}
			continue
		}

		part := formatWrapped(current)
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return err.Error()
	}

	return strings.Join(parts, ": ")
}

// Format implements fmt.Formatter for custom formatting.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprint(s, e.Details())
			return
		}
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	case 's':
		fmt.Fprint(s, e.Error())
	}
}

// formatLevel formats a single error level into a string, including its message and fields.
func formatLevel(e *Error) string {
	if e == nil {
		return ""
	}

	msg := e.msg
	fields := formatFields(e.fields)
	if fields == "" {
		return msg
	}
	if msg == "" {
		return fields
	}

	return msg + " " + fields
}

// formatFields formats the fields of an error into a string representation, sorted by key.
func formatFields(fields map[string]any) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, fields[key]))
	}

	return "(" + strings.Join(parts, ", ") + ")"
}

// formatWrapped formats a non-goerr wrapper level without duplicating its wrapped error text.
func formatWrapped(err error) string {
	if err == nil {
		return ""
	}

	next := unwrap(err)
	if next == nil {
		return err.Error()
	}

	full := err.Error()
	suffix := ": " + next.Error()
	if strings.HasSuffix(full, suffix) {
		return strings.TrimSuffix(full, suffix)
	}

	return full
}

func unwrap(err error) error {
	if err == nil {
		return nil
	}

	next, ok := err.(interface{ Unwrap() error })
	if !ok {
		return nil
	}

	return next.Unwrap()
}
