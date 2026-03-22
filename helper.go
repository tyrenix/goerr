package goerr

import (
	"errors"
)

// AsError attempts to cast err to *Error. It returns the error and true if successful, or nil and false otherwise.
func AsError(err error) (*Error, bool) {
	var target *Error
	if !errors.As(err, &target) {
		return nil, false
	}
	return target, true
}

// CodeOf returns the nearest code from the error chain.
func CodeOf(err error) (Code, bool) {
	t, ok := AsError(err)
	if !ok || t.Spec().IsZero() {
		return "", false
	}
	return t.Code(), true
}

// CodeIs reports whether the nearest code in the chain matches code.
func CodeIs(err error, code Code) bool {
	got, ok := CodeOf(err)
	return ok && got == code
}

// KindOf returns the nearest kind from the error chain.
func KindOf(err error) (Kind, bool) {
	t, ok := AsError(err)
	if !ok || t.Spec().IsZero() {
		return "", false
	}
	return t.Kind(), true
}

// KindIs reports whether the nearest kind in the chain matches kind.
func KindIs(err error, kind Kind) bool {
	got, ok := KindOf(err)
	return ok && got == kind
}

// FieldOf returns the nearest field value from the error chain.
func FieldOf(err error, key string) (any, bool) {
	t, ok := AsError(err)
	if !ok {
		return nil, false
	}

	return t.GetField(key)
}
