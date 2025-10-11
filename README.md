# goerr

`goerr` is a lightweight error wrapper for Go that extends `errors.Join`. It allows wrapping multiple errors while exposing only the primary error to users, keeping the full chain and metadata (e.g., HTTP status codes, custom fields) for logging and debugging.

## Installation

```bash
go get -u github.com/tyrenix/goerr@v1.2.0
```

## Features

- **Consistent Error Nesting**: `goerr.New(base, ...args)` now consistently nests `*goerr.Error` instances, allowing `errors.Is` to correctly identify wrapped `*goerr.Error` objects.
- **Minimal user-facing error**: `.Error()` returns only the primary error message, safe for user display.
- **Full error context**: Use `Details()` or `fmt.Printf("%v", err)` to get the full chain of errors and fields. Example output: `invalid request: validation failed: amount too low: fields: http_code=400, min_amount=100`
- **Recursive Custom fields**: Add metadata (e.g., `min_amount`, `http_code`) via `WithField` and access via `GetField` or `Fields`. These methods now recursively traverse the entire error chain.
- **HTTP status codes**: Set and retrieve HTTP status codes with `WithHTTPCode` and `HTTPCode`. `HTTPCode` also works recursively.
- **Error conversion**: Convert any error to `*goerr.Error` using `FromError`.
- **Custom formatting**: Implements `fmt.Formatter` with `%v` (full context), `%s` (primary error), `%q` (quoted primary error).
- **Go 1.20+ Compatibility**: Supports `errors.Unwrap() []error`, `errors.Is`, `errors.As` for seamless integration with modern Go error handling.

## Usage

```go
package main

import (
	"errors"
	"fmt"
	"github.com/tyrenix/goerr"
)

func main() {
	// Create an error with wrapped errors and metadata
	err := goerr.New(
		"invalid request",
		"validation failed",
		errors.New("amount too low"),
		goerr.WithField("min_amount", 100),
		goerr.WithHTTPCode(400),
	)

	// User-facing message
	fmt.Println("User message:", err.Error()) // User message: invalid request

	// Full context for logging (note: fields are now part of the recursive output)
	fmt.Printf("Full context: %v\n", err) // Full context: invalid request: validation failed: amount too low: (http_code=400, min_amount=100)

	// Convert and access fields (now works recursively)
	goErr := goerr.FromError(err)
	if minAmount, ok := goErr.GetField("min_amount"); ok {
		fmt.Println("Min amount:", minAmount) // Min amount: 100
	}
	if httpCode, ok := goErr.GetField("http_code"); ok {
		fmt.Println("HTTP status:", httpCode) // HTTP status: 400
	}

	// Check error type (now works correctly with nested errors)
	if errors.Is(err, errors.New("amount too low")) {
		fmt.Println("Found amount too low error")
	}
}
```

## API

### `goerr.New(main any, args ...any) error`
Creates an error with a primary message or error, consistently nesting `*goerr.Error` instances and joining additional strings, errors, or options.

### `goerr.FromError(err error) *Error`
Converts any error to `*goerr.Error`, preserving context if already `*goerr.Error`.

### `goerr.WithError(wrapped error) Option`
Adds a wrapped error to the error chain.

### `goerr.WithField(key string, value any) Option`
Adds a key-value pair to the error's metadata.

### `goerr.WithHTTPCode(code int) Option`
Sets the HTTP status code as a field (`http_code`).

### `(*Error).Error() string`
Returns the primary error message.

### `(*Error).Unwrap() []error`
Returns all wrapped errors (including `mainErr`) for compatibility with Go 1.20+ `errors.Is` and `errors.As`.

### `(*Error).Wrapped() []error`
Returns the list of explicitly wrapped errors (excluding `mainErr`).

### `(*Error).Fields() map[string]any`
Returns a map of all custom fields collected recursively from the entire error chain. Fields from outer errors overwrite fields from inner errors.

### `(*Error).GetField(key string) (any, bool)`
Returns the value of a field by key, searching recursively through the entire error chain. Returns a boolean indicating if it exists.

### `(*Error).HTTPCode() int`
Returns the associated HTTP status code, if set, searching recursively through the error chain.

### `(*Error).Details() string`
Returns a detailed string with the primary error, all wrapped errors, and fields, formatted recursively.

### `(*Error).Format(s fmt.State, verb rune)`
Implements `fmt.Formatter`: `%v` for full context, `%s` for primary error, `%q` for quoted primary error.

## Version
- Current: `v1.2.0`
- Backwards compatible with `v1.1.2`

## Contributing
Contributions are welcome! Please open an issue or pull request at [github.com/tyrenix/goerr](https://github.com/tyrenix/goerr).

## License
MIT — see [LICENSE](https://github.com/tyrenix/goerr/blob/master/LICENSE)
