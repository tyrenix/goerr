package main

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v3"
)

var ErrPasswordTooShort = goerr.New(
	"password is too short",
	goerr.WithSpec(goerr.Define("password.too_short", goerr.KindInvalid)),
)

func main() {
	err := validatePassword("123")
	if err == nil {
		return
	}

	code, _ := goerr.CodeOf(err)
	kind, _ := goerr.KindOf(err)

	fmt.Println("error:", err)
	fmt.Println("code:", code)
	fmt.Println("kind:", kind)
	fmt.Println("is password too short:", errors.Is(err, ErrPasswordTooShort))
}

func validatePassword(password string) error {
	if len(password) >= 6 {
		return nil
	}

	return ErrPasswordTooShort
}
