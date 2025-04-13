package goerr

// Option defines a functional option for modifying the Error struct.
type Option func(*Error)

// WithHTTPCode sets the HTTP status code for the Error.
func WithHTTPCode(code int) Option {
	return func(e *Error) {
		e.httpCode = code
	}
}
