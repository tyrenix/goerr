package goerr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorNilSafety(t *testing.T) {
	var err *Error

	require.Equal(t, "", err.Error())
	require.Equal(t, Spec{}, err.Spec())
	require.Equal(t, Code(""), err.Code())
	require.Equal(t, Kind(""), err.Kind())
	require.False(t, err.Is(New("user not found", WithSpec(Define("user.not_found", KindNotFound)))))
}

func TestErrorIsBySpec(t *testing.T) {
	spec := Define("user.not_found", KindNotFound)

	err := New("user not found", WithSpec(spec))
	target := New("another text", WithSpec(spec))

	require.True(t, errors.Is(err, target))
}

func TestNewWithSpec(t *testing.T) {
	err := NewWithSpec("user not found", "user.not_found", KindNotFound)

	goErr, ok := AsError(err)
	require.True(t, ok)
	require.Equal(t, "user not found", goErr.Error())
	require.Equal(t, Code("user.not_found"), goErr.Code())
	require.Equal(t, KindNotFound, goErr.Kind())
}

func TestErrorIsDoesNotMatchZeroSpec(t *testing.T) {
	err := New("plain error")
	target := New("another plain error")

	require.False(t, errors.Is(err, target))
}

func TestHelpersOnNestedStdlibChain(t *testing.T) {
	spec := Define("auth.unauthorized", KindUnauthorized)
	businessErr := New("unauthorized", WithSpec(spec))
	technicalErr := errors.New("token expired")

	err := fmt.Errorf("service: %w", fmt.Errorf("repo: %w: %w", businessErr, technicalErr))

	code, ok := CodeOf(err)
	require.True(t, ok)
	require.Equal(t, spec.Code, code)

	kind, ok := KindOf(err)
	require.True(t, ok)
	require.Equal(t, spec.Kind, kind)

	require.True(t, CodeIs(err, spec.Code))
	require.True(t, KindIs(err, spec.Kind))
	require.True(t, errors.Is(err, businessErr))
	require.True(t, errors.Is(err, technicalErr))
}
