package goerr

// New creates a new structured error.
func New(msg string, opts ...Option) error {
	err := &Error{msg: msg}

	for _, opt := range opts {
		opt(err)
	}

	return err
}
