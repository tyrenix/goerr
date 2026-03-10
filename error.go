package goerr

// Error is a business error with a stable code and kind.
type Error struct {
	msg  string
	spec Spec
}

// Error returns the business error message.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.msg
}

// Is reports whether target has the same non-zero specification.
func (e *Error) Is(target error) bool {
	if e == nil || target == nil || e.spec.IsZero() {
		return false
	}

	t, ok := target.(*Error)
	if !ok || t == nil || t.spec.IsZero() {
		return false
	}

	return e.spec == t.spec
}

// Spec returns the error specification.
func (e *Error) Spec() Spec {
	if e == nil {
		return Spec{}
	}

	return e.spec
}

// Code returns the error code.
func (e *Error) Code() Code {
	if e == nil {
		return ""
	}

	return e.spec.Code
}

// Kind returns the error kind.
func (e *Error) Kind() Kind {
	if e == nil {
		return ""
	}

	return e.spec.Kind
}
