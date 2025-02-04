```markdown
# goerr

The `goerr` package offers a minimalist way to join multiple errors into one using Go's standard `errors.Join`. It helps to hide detailed server errors from the user by exposing only the main error message, while the full error chain is available for logging.

## Installation

Install the package using the standard `go get` command:

```bash
go get -u github.com/tyrenix/goerr
```

## Features

- **Error Joining:**  
  The `New` function accepts a primary error and additional errors (of type `string` or `error`), joining them using `errors.Join`.

- **User-friendly Output:**  
  The `Error()` method returns only the first error message, intended for user display.

- **Logging Support:**  
  The complete error chain can be accessed via `Unwrap()`, which is useful for logging.

- **Custom Formatting:**  
  Implements the `fmt.Formatter` interface, allowing flexible formatting with `%v`, `%q`, and `%s`.

## Usage

```go
package main

import (
	"fmt"
	"goerr" // replace with your module path
)

func main() {
	// Main error to be shown to the user
	baseErr := fmt.Errorf("internal error")

	// Additional error for logs
	addErr := "database connection failed"

	// Create a joined error
	err := goerr.New(baseErr, addErr)

	// Display error to the user (only the main error)
	fmt.Println("User error:", err.Error())

	// Log the complete error chain
	fmt.Printf("Error details: %+v\n", err)
}
```

## Notes

- The package is designed with minimalism and efficiency in mind.
- The main error is displayed to the user, while additional information is kept for logging via `Unwrap()`.

## License

This project is licensed under the MIT License.
```