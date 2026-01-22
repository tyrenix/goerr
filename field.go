package goerr

import "maps"

// Fields returns the map of custom fields.
func (e *Error) Fields() map[string]any {
	// if error is nil, return nil
	if e == nil {
		return nil
	}

	// init fields
	var fields map[string]any

	// get fields from main error
	if mainErr, ok := e.mainErr.(*Error); ok {
		fields = mainErr.Fields()
	} else {
		fields = make(map[string]any)
	}

	// get fields from wrapped
	for _, w := range e.wrapped {
		if wErr, ok := w.(*Error); ok {
			maps.Copy(fields, wErr.Fields())
		}
	}

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

	// search in mainErr recursively
	if mainGoErr, ok := e.mainErr.(*Error); ok {
		if v, ok := mainGoErr.GetField(key); ok {
			return v, ok
		}
	}

	// search in wrapped errors recursively
	for _, w := range e.wrapped {
		if wrappedGoErr, ok := w.(*Error); ok {
			if v, ok := wrappedGoErr.GetField(key); ok {
				return v, ok
			}
		}
	}

	// if not found, return false
	return nil, false
}

// Kind returns the kind of the error.
func (e *Error) Kind() error {
	// if kind is nil, return kind from main error if it is an Error
	if e.kind == nil {
		if err, ok := e.mainErr.(*Error); ok {
			return err.Kind()
		}
	}
	// return kind
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
