package goerr

// Code represents a stable business error code.
type Code string

// Built-in error codes.
const (
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeInvalid      Code = "invalid"
	CodeConflict     Code = "conflict"
	CodeUnauthorized Code = "unauthorized"
	CodeInternal     Code = "internal"
)
