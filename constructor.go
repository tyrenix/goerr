package goerr

// New creates a new structured error.
func New(msg string, opts ...Option) error {
	err := &Error{
		msg:    msg,
		fields: map[string]any{},
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// Wrap wraps an existing error with a new message.
func Wrap(err error, msg string, opts ...Option) error {
	if err == nil {
		return nil
	}

	opts = append([]Option{withCause(err)}, opts...)

	return New(msg, opts...)
}
