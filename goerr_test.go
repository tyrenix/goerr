package goerr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tyrenix/goerr"
)

// customError is a simple error type used for testing errors.As
type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }

func TestGoerr(t *testing.T) {
	// --- Test Primitives ---
	baseErr := errors.New("base error")
	customErr := &customError{msg: "custom error"}

	// --- Errors for testing nesting ---
	// e1 is a simple goerr
	e1 := goerr.New("e1", goerr.WithField("f1", "v1"))

	// e2 wraps a standard error and has its own field
	e2 := goerr.New("e2", baseErr, goerr.WithField("f2", "v2"), goerr.WithHTTPCode(404))

	// e3 wraps e1, creating a nested goerr
	e3 := goerr.New("e3", e1, goerr.WithField("f3", "v3"))

	// init test cases
	tests := []struct {
		name string
		err  error

		// Expectations
		expectMsg      string
		expectDetails  string // for %v
		expectIs       []error
		expectAs       any
		expectFields   map[string]any
		expectHTTPCode int
	}{
		{
			name:          "Simple string error",
			err:           goerr.New("simple error"),
			expectMsg:     "simple error",
			expectDetails: "simple error",
		},
		{
			name:          "Error with wrapped string and standard error",
			err:           goerr.New("main", "wrapped string", baseErr),
			expectMsg:     "main",
			expectDetails: "main: wrapped string: base error",
			expectIs:      []error{baseErr},
		},
		{
			name:           "Error with options",
			err:            goerr.New("with options", goerr.WithField("key", "value"), goerr.WithHTTPCode(400)),
			expectMsg:      "with options",
			expectDetails:  "with options: (http_code=400, key=value)",
			expectFields:   map[string]any{"key": "value", "http_code": 400},
			expectHTTPCode: 400,
		},
		{
			name:          "errors.As check",
			err:           goerr.New("main", customErr),
			expectMsg:     "main",
			expectDetails: "main: custom error",
			expectAs:      &customError{},
		},
		{
			name:          "Nested goerr as argument",
			err:           e3, // e3 wraps e1
			expectMsg:     "e3",
			expectDetails: "e3: e1: (f1=v1): (f3=v3)",
			expectIs:      []error{e1},
			expectFields:  map[string]any{"f1": "v1", "f3": "v3"},
		},
		{
			name:           "Nested goerr as main error",
			err:            goerr.New(e2, "outer layer"), // e2 is the main error
			expectMsg:      "e2",
			expectDetails:  "e2: base error: (f2=v2, http_code=404): outer layer",
			expectIs:       []error{e2, baseErr},
			expectFields:   map[string]any{"f2": "v2", "http_code": 404},
			expectHTTPCode: 404,
		},
		{
			name:         "GetField recursive from mainErr",
			err:          goerr.New(e2, "outer"), // e2 has f2="v2"
			expectFields: map[string]any{"f2": "v2"},
		},
		{
			name:         "GetField recursive from wrapped err",
			err:          goerr.New("outer", e2), // e2 has f2="v2"
			expectFields: map[string]any{"f2": "v2"},
		},
	}

	// run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// if err is nil, fail
			if tc.err == nil {
				t.Fatal("error is nil")
			}

			// test .Error() method
			if tc.expectMsg != "" && tc.err.Error() != tc.expectMsg {
				t.Errorf("Expected .Error() to be %q, but got %q", tc.expectMsg, tc.err.Error())
			}

			// test formatting (%v)
			details := fmt.Sprintf("%v", tc.err)
			if tc.expectDetails != "" && details != tc.expectDetails {
				t.Errorf("Expected details to be %q, but got %q", tc.expectDetails, details)
			}

			// test errors.Is
			for _, target := range tc.expectIs {
				if !errors.Is(tc.err, target) {
					t.Errorf("Expected errors.Is to be true for target %q", target)
				}
			}

			// test errors.As
			if tc.expectAs != nil {
				if !errors.As(tc.err, &tc.expectAs) {
					t.Errorf("Expected errors.As to be true for target type %T", tc.expectAs)
				}
			}

			// test field methods
			if goErr, ok := tc.err.(*goerr.Error); ok {
				// test HTTPCode()
				if tc.expectHTTPCode != 0 && goErr.HTTPCode() != tc.expectHTTPCode {
					t.Errorf("Expected HTTPCode() to be %d, but got %d", tc.expectHTTPCode, goErr.HTTPCode())
				}

				// test GetField() and Fields()
				if tc.expectFields != nil {
					allFields := goErr.Fields()
					for key, expectedValue := range tc.expectFields {
						// test GetField
						value, ok := goErr.GetField(key)
						if !ok {
							t.Errorf("GetField: expected to find key %q, but did not", key)
							continue
						}
						if value != expectedValue {
							t.Errorf("GetField: expected value for key %q to be %v, but got %v", key, expectedValue, value)
						}

						// Test Fields
						value, ok = allFields[key]
						if !ok {
							t.Errorf("Fields: expected to find key %q, but did not", key)
							continue
						}
						if value != expectedValue {
							t.Errorf("Fields: expected value for key %q to be %v, but got %v", key, expectedValue, value)
						}
					}
				}
			}
		})
	}
}