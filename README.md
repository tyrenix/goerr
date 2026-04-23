# goerr

`goerr` is a small business error package for Go.

It does not replace `fmt.Errorf`. Use `fmt.Errorf` for context and wrapping, and use `goerr` to attach stable business classification through `Code` and `Kind`.

## Installation

```bash
go get github.com/daanila01/goerr@latest
```

## API

```go
type Code string
type Kind string

type Spec struct {
	Code Code
	Kind Kind
}

func New(msg string, opts ...Option) error
func NewWithSpec(msg string, code Code, kind Kind, opts ...Option) error
func Define(code Code, kind Kind) Spec
func WithSpec(spec Spec) Option

func AsError(err error) (*Error, bool)
func CodeOf(err error) (Code, bool)
func CodeIs(err error, code Code) bool
func KindOf(err error) (Kind, bool)
func KindIs(err error, kind Kind) bool
```

## Usage

Define business errors in your shared or domain packages:

```go
var (
	ErrInternal = goerr.NewWithSpec("internal", "internal", goerr.KindInternal)
	ErrUserNotFound = goerr.NewWithSpec("user not found", "user.not_found", goerr.KindNotFound)
)
```

`NewWithSpec` is the preferred shorthand for predefined business errors. Use `New(..., WithSpec(...))` when you already have a reusable `Spec`.

Map technical errors to business errors at boundaries:

```go
func (r *Repo) GetByID(id string) error {
	err := queryUser(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("execute user query: %w: %w", ErrUserNotFound, err)
		}

		return fmt.Errorf("execute user query: %w: %w", ErrInternal, err)
	}

	return nil
}
```

Add service context with plain `fmt.Errorf`:

```go
func (s *Service) GetByID(id string) error {
	if err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	return nil
}
```

Extract business data on the transport layer:

```go
func writeError(err error) {
	code, _ := goerr.CodeOf(err)
	kind, _ := goerr.KindOf(err)

	fmt.Println("code:", code)
	fmt.Println("kind:", kind)
}
```

## Notes

- `Error()` returns the developer-facing message of the matched business error.
- `CodeOf` and `KindOf` return the nearest `goerr.Error` found by `errors.As`.
- `CodeIs` and `KindIs` are helpers over `CodeOf` and `KindOf`.
- `errors.Is` still works for both the business error and the technical cause when they are wrapped with `%w`.
- if you rely on the nearest business error, keep it first in the `%w` chain
