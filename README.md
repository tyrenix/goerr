# goerr

`goerr` is a lightweight error wrapper for Go that builds on top of `errors.Join`. It provides a clean way to wrap multiple errors while exposing only the primary error to users and keeping the full chain available for logging and debugging. It also supports optional metadata like HTTP status codes.

## Installation

```bash
go get -u github.com/tyrenix/goerr
```

## Features

- **Error joining:**  
  `goerr.New(base, ...opts)` joins a base error with additional messages or wrapped errors.

- **Minimal user-facing error:**  
  The `.Error()` method returns only the first (main) message — safe for user display.

- **Full error chain for logs:**  
  Use `fmt.Printf("%+v", err)` or `Unwrap()` to get the full context.

- **Optional metadata (e.g., HTTP codes):**  
  Pass `WithHTTPCode(int)` to associate HTTP status with the error.

- **Custom formatting:**  
  Implements `fmt.Formatter`: use `%v`, `%s`, `%q` as needed.

## Usage

```go
package main

import (
	"fmt"
	"github.com/tyrenix/goerr"
)

func main() {
	base := fmt.Errorf("something went wrong")

	// Join multiple layers and add metadata
	err := goerr.New(
		base,
		"db timeout",
		fmt.Errorf("connection refused"),
		goerr.WithHTTPCode(500),
	)

	fmt.Println("User message:", err.Error())
	fmt.Printf("Full chain: %+v\n", err)

	if ge, ok := err.(*goerr.Error); ok {
		fmt.Println("HTTP status:", ge.HTTPCode())
	}
}
```

## API

### `goerr.New(base error, opts ...any) error`
Joins `base` with additional strings/errors and applies any `Option`.

### `goerr.WithHTTPCode(code int) Option`
Adds an HTTP status code to the error.

### `(*Error).Error() string`
Returns the top-level error only (first line).

### `(*Error).Unwrap() error`
Returns the full joined error for use with `errors.Unwrap()` and `errors.Is()`.

### `(*Error).HTTPCode() int`
Returns the associated HTTP status code, if set.

## License

MIT — see [LICENSE](https://github.com/tyrenix/goerr/blob/master/LICENSE)