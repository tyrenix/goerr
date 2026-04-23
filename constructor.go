package goerr

// New creates a new business error.
func New(msg string, opts ...Option) error {
	err := &Error{msg: msg}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// NewWithSpec creates a new business error with the given spec.
func NewWithSpec(msg string, code Code, kind Kind, opts ...Option) error {
	return New(msg, append(opts, WithSpec(Define(code, kind)))...)
}
