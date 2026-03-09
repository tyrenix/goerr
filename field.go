package goerr

import "maps"

// Fields returns the fields of the current error level.
func (e *Error) Fields() map[string]any {
	if e == nil {
		return nil
	}

	fields := make(map[string]any, len(e.fields))
	maps.Copy(fields, e.fields)

	return fields
}

// Field returns a field value from the current error level.
func (e *Error) Field(key string) (any, bool) {
	if e == nil {
		return nil, false
	}

	if v, ok := e.fields[key]; ok {
		return v, ok
	}

	return nil, false
}
