package goerr

import "maps"

// Fields returns the map of custom fields.
func (e *Error) Fields() map[string]any {
	// if error is nil, return nil
	if e == nil {
		return nil
	}

	// init fields
	fields := make(map[string]any, len(e.fields))
	// copy fields
	maps.Copy(fields, e.fields)

	// return fields
	return fields
}

// GetField returns the value of a field by key, with a boolean indicating if it exists.
// It searches for the key recursively through the entire error chain.
func (e *Error) GetField(key string) (any, bool) {
	// if error is nil, return nil
	if e == nil {
		return nil, false
	}

	// search in current error's fields (highest priority)
	if v, ok := e.fields[key]; ok {
		return v, ok
	}

	// return nil if not found
	return nil, false
}

// Kind returns the kind of the error.
func (e *Error) Kind() error {
	return e.kind
}

// Deprecated: HTTPCode is deprecated, this method is deprecated and will be removed in a future version.
func (e *Error) HTTPCode() int {
	if v, ok := e.GetField(fieldHTTPCode); ok {
		if httpCode, ok := v.(int); ok {
			return httpCode
		}
	}
	return 0
}
