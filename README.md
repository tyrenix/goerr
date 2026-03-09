# goerr

`goerr` is a small structured error package for Go.

It keeps the usual `New` and `Wrap` flow, but adds stable business classification through `Spec`, `Code` and `Kind`.

## Goals

- stay close to regular Go error wrapping
- classify technical errors into business codes at boundaries
- attach operation names and fields for logs
- keep public code separate from `Error()`

## Installation

```bash
go get github.com/tyrenix/goerr/v2@latest
```

## Core Types

```go
type Code string
type Kind string

type Spec struct {
	Code Code
	Kind Kind
}
```

Define business specs in your domain packages:

```go
var (
	UserNotFound = goerr.Define("user.not_found", goerr.KindNotFound)
	UserExists   = goerr.Define("user.exists", goerr.KindConflict)
	Internal     = goerr.Define("internal", goerr.KindInternal)
)
```

Import path:

```go
import "github.com/tyrenix/goerr/v2"
```

## API

```go
func New(msg string, opts ...Option) error
func Wrap(err error, msg string, opts ...Option) error

func Define(code Code, kind Kind) Spec

func WithSpec(spec Spec) Option
func WithOp(op string) Option
func WithField(key string, value any) Option
func WithFields(fields map[string]any) Option

func CodeOf(err error) (Code, bool)
func KindOf(err error) (Kind, bool)
func FieldOf(err error, key string) (any, bool)
func AllFields(err error) map[string]any
func AsError(err error) (*Error, bool)
func DetailsOf(err error) string
```

## Usage

### Validation

Use `New` when you are creating an error directly without a cause.

```go
var PasswordTooShort = goerr.Define("password.too_short", goerr.KindInvalid)

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return goerr.New(
			"password is too short",
			goerr.WithSpec(PasswordTooShort),
			goerr.WithOp("auth.ValidatePassword"),
		)
	}

	return nil
}
```

### Repository Boundary

Assign business classification where a technical error first becomes a domain error.

```go
var UserNotFound = goerr.Define("user.not_found", goerr.KindNotFound)
var Internal = goerr.Define("internal", goerr.KindInternal)

func (r *Repo) GetUser(ctx context.Context, id string) (*User, error) {
	user, err := r.db.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, goerr.Wrap(
				err,
				"execute user query",
				goerr.WithSpec(UserNotFound),
				goerr.WithOp("userrepo.GetUser"),
				goerr.WithField("user_id", id),
			)
		}

		return nil, goerr.Wrap(
			err,
			"execute user query",
			goerr.WithSpec(Internal),
			goerr.WithOp("userrepo.GetUser"),
			goerr.WithField("user_id", id),
		)
	}

	return user, nil
}
```

### Service Layer

If an error is already classified below, just add context.

```go
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, goerr.Wrap(
			err,
			"get user",
			goerr.WithOp("userservice.GetUser"),
		)
	}

	return user, nil
}
```

### Handler Layer

Use `CodeOf` and `KindOf` for public output and transport mapping.

```go
func writeError(w http.ResponseWriter, err error) {
	code, _ := goerr.CodeOf(err)
	kind, _ := goerr.KindOf(err)

	status := http.StatusInternalServerError
	switch kind {
	case goerr.KindInvalid:
		status = http.StatusBadRequest
	case goerr.KindNotFound:
		status = http.StatusNotFound
	}

	http.Error(w, string(code), status)
}
```

## Output

```go
err := goerr.Wrap(
	sql.ErrNoRows,
	"execute user query",
	goerr.WithSpec(goerr.Define("user.not_found", goerr.KindNotFound)),
	goerr.WithOp("userrepo.GetUser"),
)

err = goerr.Wrap(err, "get user", goerr.WithOp("userservice.GetUser"))
```

`err.Error()`:

```text
get user: execute user query: sql: no rows in result set
```

`fmt.Printf("%+v\n", err)`:

```text
user.not_found (kind=not_found): get user (op=userservice.GetUser): execute user query (op=userrepo.GetUser): sql: no rows in result set
```

## Notes

- `Wrap` returns `nil` when the incoming error is `nil`
- `errors.Is` and `errors.As` still work with the wrapped technical cause
- `Error()` is developer-facing
- `%v` uses `Error()`
- `%+v` uses detailed structured output for `goerr` values
- `DetailsOf(err)` is the safest way to get detailed output from mixed error chains
- `CodeOf` is the stable business code for localization and public responses
