package goerr

// Error is a structured error with message, specification and cause.
type Error struct {
	msg    string
	spec   Spec
	cause  error
	fields map[string]any
}

// Error returns the error message chain.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.msg == "" {
		if e.cause == nil {
			return ""
		}
		return e.cause.Error()
	}

	if e.cause == nil {
		return e.msg
	}

	return e.msg + ": " + e.cause.Error()
}

// Unwrap returns the wrapped cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}

// Message returns the message at the current level.
func (e *Error) Message() string {
	if e == nil {
		return ""
	}

	return e.msg
}

// Spec returns the current level specification.
func (e *Error) Spec() Spec {
	if e == nil {
		return Spec{}
	}

	return e.spec
}

// Code returns the current level code.
func (e *Error) Code() Code {
	if e == nil {
		return ""
	}

	return e.spec.Code
}

// Kind returns the current level kind.
func (e *Error) Kind() Kind {
	if e == nil {
		return ""
	}

	return e.spec.Kind
}
