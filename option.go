package goerr

// Option defines a functional option for modifying the Error struct.
type Option func(*Error)

// WithSpec sets the business specification of the error.
func WithSpec(spec Spec) Option {
	return func(e *Error) {
		e.spec = spec
	}
}

// WithField attaches a single public field to the error.
func WithField(key string, value any) Option {
	return func(e *Error) {
		if key == "" {
			return
		}

		if e.fields == nil {
			e.fields = make(map[string]any, 1)
		}

		e.fields[key] = value
	}
}

// WithFields attaches multiple public fields to the error.
func WithFields(fields map[string]any) Option {
	return func(e *Error) {
		if len(fields) == 0 {
			return
		}

		if e.fields == nil {
			e.fields = make(map[string]any, len(fields))
		}

		for key, value := range fields {
			if key == "" {
				continue
			}

			e.fields[key] = value
		}
	}
}
