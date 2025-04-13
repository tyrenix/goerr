package goerr_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/tyrenix/goerr"
)

func TestNew_WithStringAndError(t *testing.T) {
	base := errors.New("base")
	wrapped := errors.New("wrapped")
	err := goerr.New(base, "additional", wrapped)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// check string
	if !strings.Contains(err.Error(), "base") {
		t.Error("expected 'base' in error")
	}
}

func TestNew_WithHTTPCode(t *testing.T) {
	base := errors.New("base")
	err := goerr.New(base, goerr.WithHTTPCode(403))

	e, ok := err.(*goerr.Error)
	if !ok {
		t.Fatalf("expected *goerr.Error, got %T", err)
	}

	if e.HTTPCode() != 403 {
		t.Errorf("expected httpCode 403, got %d", e.HTTPCode())
	}
}

func TestUnwrap(t *testing.T) {
	base := errors.New("base")
	err := goerr.New(base)

	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		t.Error("expected unwrapped error, got nil")
	}
	if !strings.Contains(unwrapped.Error(), "base") {
		t.Error("unwrapped error missing base")
	}
}

func TestFormat(t *testing.T) {
	base := errors.New("base")
	err := goerr.New(base, "layer 1", "layer 2")

	formatted := fmt.Sprintf("%v", err)
	if !strings.Contains(formatted, "layer 1") || !strings.Contains(formatted, "layer 2") {
		t.Errorf("unexpected format: %s", formatted)
	}

	quoted := fmt.Sprintf("%q", err)
	if !strings.Contains(quoted, "\"base\"") {
		t.Errorf("unexpected quoted format: %s", quoted)
	}
}
