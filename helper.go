package goerr

import (
	"errors"
	"maps"
)

// AsError returns the first goerr.Error from the chain.
func AsError(err error) (*Error, bool) {
	var target *Error
	if !errors.As(err, &target) {
		return nil, false
	}

	return target, true
}

// CodeOf returns the nearest code from the error chain.
func CodeOf(err error) (Code, bool) {
	for _, item := range chain(err) {
		if item.spec.Code != "" {
			return item.spec.Code, true
		}
	}

	return "", false
}

// CodeIs reports whether the nearest code in the chain matches code.
func CodeIs(err error, code Code) bool {
	got, ok := CodeOf(err)
	return ok && got == code
}

// KindOf returns the nearest kind from the error chain.
func KindOf(err error) (Kind, bool) {
	for _, item := range chain(err) {
		if item.spec.Kind != "" {
			return item.spec.Kind, true
		}
	}

	return "", false
}

// KindIs reports whether the nearest kind in the chain matches kind.
func KindIs(err error, kind Kind) bool {
	got, ok := KindOf(err)
	return ok && got == kind
}

// FieldOf returns the nearest field value from the error chain.
func FieldOf(err error, key string) (any, bool) {
	for _, item := range chain(err) {
		if v, ok := item.Field(key); ok {
			return v, true
		}
	}

	return nil, false
}

// AllFields returns merged fields from the full error chain.
func AllFields(err error) map[string]any {
	items := chain(err)
	if len(items) == 0 {
		return nil
	}

	fields := map[string]any{}
	for i := len(items) - 1; i >= 0; i-- {
		maps.Copy(fields, items[i].fields)
	}

	return fields
}

// chain returns the error chain as a slice of *Error, starting from the outermost error.
func chain(err error) []*Error {
	var items []*Error
	for err != nil {
		if current, ok := err.(*Error); ok {
			items = append(items, current)
		}

		next, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}

		err = next.Unwrap()
	}

	return items
}
