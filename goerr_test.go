package goerr_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/tyrenix/goerr"
)

// customError is a simple error type used for testing errors.As
type customError struct {
	msg string
}

func (e *customError) Error() string { return e.msg }

func TestErrorsIsAndAs(t *testing.T) {
	// Case 1: using raw errors.New as base
	base := errors.New("base error")
	baseGoErr := goerr.New(base)
	custom := &customError{"custom wrapped error"}

	baseErr := goerr.New(base, custom, "extra info")
	goErr := goerr.New(baseGoErr, custom, "extra info")

	// Test errors.Is with original base
	if !errors.Is(baseErr, base) {
		t.Errorf("errors.Is failed: expected to find base error in chain")
	}
	// Test errors.Is with goerr
	if !errors.Is(goErr, baseGoErr) {
		t.Errorf("errors.Is failed: expected to find goerr error in chain")
	}

	// Test errors.As with custom error
	var target *customError
	if !errors.As(baseErr, &target) {
		t.Errorf("errors.As failed: expected to find customError in chain")
	}
	if target.msg != "custom wrapped error" {
		t.Errorf("unexpected custom error message: got %q", target.msg)
	}

	// Case 2: using goerr.New as base
	e1 := goerr.New("main error")
	e2 := goerr.New(e1, fmt.Errorf("nested db failure"))

	// Test errors.Is with goerr wrapped error
	if !errors.Is(e2, e1) {
		t.Errorf("errors.Is failed: expected to find goerr error in chain")
	}

	// Test that the main error message is preserved
	if e1.Error() != "main error" {
		t.Errorf("unexpected main error message: %q", e1.Error())
	}
}

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
