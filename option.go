package goerr

import "maps"

// Option defines a functional option for modifying the Error struct.
type Option func(*Error)

// fieldHTTPCode is the key for the HTTP status code field.
const fieldHTTPCode = "http_code"

// Deprecated: WithError is deprecated, this method is deprecated and will be removed in a future version.
func WithError(wrapped error) Option {
	return func(e *Error) {
		if wrapped != nil {
			e.wrapped = wrapped
		}
	}
}

// Kind sets the kind of the Error.
func Kind(kind error) Option {
	return func(e *Error) {
		e.kind = kind
	}
}

// Deprecated: WithField is deprecated, use Field instead.
func WithField(key string, value any) Option {
	return Field(key, value)
}

// Fields adds multiple key-value pairs to the Error's fields.
func Fields(fields map[string]any) Option {
	return func(e *Error) {
		maps.Copy(e.fields, fields)
	}
}

// Field adds a key-value pair to the Error's fields.
func Field(key string, value any) Option {
	return func(e *Error) {
		e.fields[key] = value
	}
}

// Deprecated: WithHTTPCode is deprecated, this method is deprecated and will be removed in a future version.
func WithHTTPCode(code int) Option {
	return WithField(fieldHTTPCode, code)
}
