# goerr

`goerr` is a lightweight error wrapper for Go that extends `errors.Join`. It allows wrapping multiple errors while exposing only the primary error to users, keeping the full chain and metadata (e.g., HTTP status codes, custom fields) for logging and debugging.

## Installation

```bash
go get -u github.com/tyrenix/goerr@v1.2.0
```

## Features

- **Error joining**: Combine a primary error with additional errors or messages using `goerr.New(base, ...args)`.
- **Minimal user-facing error**: `.Error()` returns only the primary error message, safe for user display.
- **Full error context**: Use `Details()` or `fmt.Printf("%v", err)` to get the full chain of errors and fields.
- **Custom fields**: Add metadata (e.g., `min_amount`, `http_code`) via `WithField` and access via `GetField` or `Fields`.
- **HTTP status codes**: Set and retrieve HTTP status codes with `WithHTTPCode` and `HTTPCode`.
- **Error conversion**: Convert any error to `*goerr.Error` using `FromError`.
- **Custom formatting**: Implements `fmt.Formatter` with `%v` (full context), `%s` (primary error), `%q` (quoted primary error).
- **Go compatibility**: Supports `errors.Unwrap`, `errors.Is`, `errors.As` for seamless integration.

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

	// Full context for logging
	fmt.Printf("Full context: %v\n", err) // Full context: invalid request; validation failed; amount too low; fields: http_code=400, min_amount=100

	// Convert and access fields
	goErr := goerr.FromError(err)
	if minAmount, ok := goErr.GetField("min_amount"); ok {
		fmt.Println("Min amount:", minAmount) // Min amount: 100
	}
	if httpCode, ok := goErr.GetField("http_code"); ok {
		fmt.Println("HTTP status:", httpCode) // HTTP status: 400
	}

	// Access wrapped errors
	for _, w := range goErr.Wrapped() {
		fmt.Println("Wrapped:", w) // Wrapped: validation failed, Wrapped: amount too low
	}

	// Check error type
	if errors.Is(err, errors.New("amount too low")) {
		fmt.Println("Found amount too low error")
	}
}
```

## API

### `goerr.New(main any, args ...any) error`
Creates an error with a primary message or error, joining additional strings, errors, or options. Preserves context from existing `*goerr.Error`.

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

### `(*Error).Unwrap() error`
Returns the primary error for use with `errors.Unwrap`.

### `(*Error).Wrapped() []error`
Returns the list of wrapped errors.

### `(*Error).Fields() map[string]any`
Returns the map of custom fields.

### `(*Error).GetField(key string) (any, bool)`
Returns the value of a field by key, with a boolean indicating if it exists.

### `(*Error).HTTPCode() int`
Returns the associated HTTP status code, if set.

### `(*Error).Details() string`
Returns a detailed string with the primary error, wrapped errors, and fields.

### `(*Error).Format(s fmt.State, verb rune)`
Implements `fmt.Formatter`: `%v` for full context, `%s` for primary error, `%q` for quoted primary error.

### `(*Error).Is(target error) bool`
Supports `errors.Is` to check if the error or wrapped errors match a target.

### `(*Error).As(target any) bool`
Supports `errors.As` to assign the error or wrapped errors to a target type.

## Version
- Current: `v1.2.0`
- Backwards compatible with `v1.1.2`

## Contributing
Contributions are welcome! Please open an issue or pull request at [github.com/tyrenix/goerr](https://github.com/tyrenix/goerr).

## License
MIT — see [LICENSE](https://github.com/tyrenix/goerr/blob/master/LICENSE)