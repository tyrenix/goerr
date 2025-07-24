package goerr

// Option defines a functional option for modifying the Error struct.
type Option func(*Error)

// fieldHTTPCode is the key for the HTTP status code field.
const fieldHTTPCode = "http_code"

// WithError adds a wrapped error to the Error.
func WithError(wrapped error) Option {
	return func(e *Error) {
		if wrapped != nil {
			e.wrapped = append(e.wrapped, wrapped)
		}
	}
}

// WithField adds a key-value pair to the Error's fields.
func WithField(key string, value any) Option {
	return func(e *Error) {
		e.fields[key] = value
	}
}

// WithHTTPCode sets the HTTP status code for the Error.
func WithHTTPCode(code int) Option {
	return WithField("http_code", code)
}
