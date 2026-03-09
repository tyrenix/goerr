package goerr

// Kind represents an error kind.
type Kind string

// Built-in error kinds.
const (
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindInvalid      Kind = "invalid"
	KindConflict     Kind = "conflict"
	KindUnauthorized Kind = "unauthorized"
	KindInternal     Kind = "internal"
)
