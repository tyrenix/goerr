package goerr

// Spec describes the public classification of an error.
type Spec struct {
	Code Code
	Kind Kind
}

// Define creates a new error specification.
func Define(code Code, kind Kind) Spec {
	return Spec{
		Code: code,
		Kind: kind,
	}
}

// IsZero reports whether the specification is empty.
func (s Spec) IsZero() bool {
	return s.Code == "" && s.Kind == ""
}
