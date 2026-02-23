package goerr

// Built-in sentinel errors for cross-package errors.Is checks.
var (
	ErrForbidden    = New("forbidden", Kind(KindForbidden))
	ErrNotFound     = New("not_found", Kind(KindNotFound))
	ErrInvalid      = New("invalid", Kind(KindInvalid))
	ErrConflict     = New("conflict", Kind(KindConflict))
	ErrUnauthorized = New("unauthorized", Kind(KindUnauthorized))
	ErrInternal     = New("internal", Kind(KindInternal))
)
