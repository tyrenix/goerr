package goerr

// KindValue represents an error kind
type KindValue string

// Built-in error kinds.
const (
	KindForbidden    KindValue = "forbidden"
	KindNotFound     KindValue = "not_found"
	KindInvalid      KindValue = "invalid"
	KindConflict     KindValue = "conflict"
	KindUnauthorized KindValue = "unauthorized"
	KindInternal     KindValue = "internal"
)
