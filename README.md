# goerr

`goerr` is a lightweight, structured error handling package for Go that provides error wrapping with context, metadata, and error kind classification for multi-transport applications (HTTP, gRPC, etc.).

## Installation

```bash
go get -u github.com/tyrenix/goerr@latest
```

## Features

- **Error Kind Classification**: Categorize errors with semantic kinds (e.g., `KindNotFound`, `KindInternal`) for consistent mapping across different transports (HTTP, gRPC, WebSocket).
- **Contextual Error Wrapping**: Add context at each layer while preserving the original error kind and metadata.
- **Structured Metadata**: Attach custom fields to errors for rich logging and debugging.
- **Minimal User-Facing Messages**: `.Error()` returns only the primary error code, safe for client display.
- **Detailed Logging**: Use `Details()` or `fmt.Printf("%v", err)` to get the full error chain with context and fields.
- **Level-Scoped Metadata**: Fields belong to the current error level; no implicit recursive lookup.
- **Go Compatible**: Full support for `errors.Is`, `errors.As` via standard `Unwrap()` chaining.
- **Custom Formatting**: Implements `fmt.Formatter` with `%v` (full details), `%s` (primary error), `%q` (quoted).

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/tyrenix/goerr"
)

// Define error kinds for your application
var (
    KindNotFound = errors.New("not_found")
    KindInternal = errors.New("internal")
)

// Define domain errors
var (
    ErrUserNotFound = goerr.New("user.not_found", goerr.Kind(KindNotFound))
    ErrInternal     = goerr.New("internal", goerr.Kind(KindInternal))
)

func main() {
    // Repository returns a domain error
    err := ErrUserNotFound

    // Service adds context with fields
    err = goerr.Wrap(err, "failed to get user",
        goerr.Field("user_id", 123),
        goerr.Field("action", "login"))

    // Handler adds more context
    err = goerr.Wrap(err, "GET /api/users/:id failed",
        goerr.Field("endpoint", "/api/users/123"),
        goerr.Field("method", "GET"))

    // Client sees only the primary business error
    fmt.Println("Error:", err.Error())
    // Output: user.not_found

    // Logs show full context with fields
    fmt.Printf("Details: %v\n", err)
    // Output: user.not_found (kind=not_found): GET /api/users/:id failed (endpoint=/api/users/123, method=GET): failed to get user (action=login, user_id=123)

    // Access error kind for transport mapping
    goErr := goerr.FromError(err)
    fmt.Println("Kind:", goErr.Kind())
    // Output: not_found

    // Access specific fields
    if userID, ok := goErr.GetField("user_id"); ok {
        fmt.Println("User ID:", userID)
        // Output: User ID: 123
    }
}
```

## Usage Patterns

### Creating Domain Errors

```go
// Define error kinds
var (
    KindNotFound = errors.New("not_found")
    KindInvalid  = errors.New("invalid")
    KindInternal = errors.New("internal")
)

// Define domain-specific errors with kinds
var (
    ErrUserNotFound  = goerr.New("user.not_found", goerr.Kind(KindNotFound))
    ErrUserBanned    = goerr.New("user.banned", goerr.Kind(KindInvalid))
    ErrOrderExpired  = goerr.New("order.expired", goerr.Kind(KindInvalid))
    ErrDatabaseError = goerr.New("internal", goerr.Kind(KindInternal))
)
```

### Wrapping Errors with Context

```go
// Repository layer
func (r *Repository) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := r.db.Find(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, ErrDatabaseError
    }
    return user, nil
}

// Service layer
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.GetUser(ctx, id)
    if err != nil {
        return nil, goerr.Wrap(err, "failed to get user",
            goerr.Field("user_id", id))
    }
    return user, nil
}

// Handler layer
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    user, err := h.service.GetUser(r.Context(), id)
    if err != nil {
        log.Error("request failed",
            "error", err, // Full context in logs
            "request_id", middleware.GetReqID(r.Context()))

        // Map error kind to HTTP status
        statusCode := MapErrorToHTTP(err)
        http.Error(w, err.Error(), statusCode)
        return
    }
    // ...
}
```

### Mapping Error Kinds to Transport Codes

```go
// HTTP mapping
func MapErrorToHTTP(err error) int {
    goErr := goerr.FromError(err)
    switch goErr.Kind() {
    case KindNotFound:
        return http.StatusNotFound
    case KindInvalid:
        return http.StatusBadRequest
    case KindInternal:
        return http.StatusInternalServerError
    default:
        return http.StatusInternalServerError
    }
}

// gRPC mapping
func MapErrorToGRPC(err error) codes.Code {
    goErr := goerr.FromError(err)
    switch goErr.Kind() {
    case KindNotFound:
        return codes.NotFound
    case KindInvalid:
        return codes.InvalidArgument
    case KindInternal:
        return codes.Internal
    default:
        return codes.Unknown
    }
}
```

## API Reference

### Creating Errors

#### `goerr.New(main any, args ...any) error`

Creates a new error with a primary message/error and optional wrapped errors or options.

```go
err := goerr.New("user.not_found",
    goerr.Kind(KindNotFound),
    goerr.Field("user_id", 123))
```

#### `goerr.Wrap(main any, context any, opts ...Option) error`

Wraps an error with additional context and fields. Preserves the original error's kind.

```go
err := goerr.Wrap(dbErr, "failed to query database",
    goerr.Field("query", "SELECT * FROM users"),
    goerr.Field("duration_ms", 150),
)
```

#### `goerr.FromError(err error) *Error`

Converts any error to `*goerr.Error`, preserving context if already a goerr error.

### Options

#### `goerr.Kind(kind error) Option`

Sets the error kind for transport mapping.

#### `goerr.Field(key string, value any) Option`

Adds a key-value metadata field.

#### `goerr.Fields(fields map[string]any) Option`

Adds multiple metadata fields at once.

```go
err := goerr.New("error",
    goerr.Fields(map[string]any{
        "user_id": 123,
        "action": "login",
        "ip": "192.168.1.1",
    }),
)
```

### Methods

#### `(*Error).Error() string`

Returns the primary business error message (safe for client display).

#### `(*Error).Details() string`

Returns the full error chain with context and fields, intended for logs and diagnostics.

#### `(*Error).Kind() error`

Returns the error kind for transport mapping.

#### `(*Error).GetField(key string) (any, bool)`

Retrieves a field value from the current error level only.

#### `(*Error).Fields() map[string]any`

Returns a copy of fields attached to the current error level.

#### `(*Error).Unwrap() error`

Returns the next wrapped error in the chain.

#### `(*Error).Format(s fmt.State, verb rune)`

Custom formatter: `%v` for details, `%s` for primary error, `%q` for quoted.

### Deprecated Methods

The following methods are deprecated and will be removed in a future version:

- `WithError(error) Option` - use wrapping via `New()` or `Wrap()` instead
- `WithField(key, value) Option` - use `Field()` instead
- `WithHTTPCode(int) Option` - use `Kind()` with error kind mapping instead
- `(*Error).HTTPCode() int` - use `Kind()` and transport-specific mapping instead

## Why goerr?

Traditional Go error handling with `fmt.Errorf` and `%w` has limitations:

- **No structured metadata**: Can't attach fields like `user_id` or `request_id`
- **Verbose error chains**: User-facing messages are intentionally minimal.
- **No error classification**: Hard to map errors to different transport codes (HTTP, gRPC)

Full context is available only through formatting (`%v`) or `Details()`.

`goerr` solves these problems by:

- Separating user-facing error codes from internal context
- Providing structured metadata for rich logging
- Classifying errors by kind for consistent transport mapping
- Maintaining full compatibility with standard Go error handling

## Contributing

Contributions are welcome! Please open an issue or pull request at [github.com/tyrenix/goerr](https://github.com/tyrenix/goerr).

## License

MIT — see [LICENSE](https://github.com/tyrenix/goerr/blob/master/LICENSE)
