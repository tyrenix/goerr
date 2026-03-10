package goerr

// Option defines a functional option for modifying the Error struct.
type Option func(*Error)

// WithSpec sets the business specification of the error.
func WithSpec(spec Spec) Option {
	return func(e *Error) {
		e.spec = spec
	}
}
