package main

import (
	"fmt"

	"github.com/tyrenix/goerr/v2"
)

var PasswordTooShort = goerr.Define("password.too_short", goerr.KindInvalid)

func main() {
	err := validatePassword("123")
	if err == nil {
		return
	}

	fmt.Println("error:", err.Error())

	if code, ok := goerr.CodeOf(err); ok {
		fmt.Println("code:", code)
	}

	if kind, ok := goerr.KindOf(err); ok {
		fmt.Println("kind:", kind)
	}

	fmt.Printf("details: %+v\n", err)
}

func validatePassword(password string) error {
	if len(password) >= 6 {
		return nil
	}

	return goerr.New(
		"password is too short",
		goerr.WithSpec(PasswordTooShort),
		goerr.WithOp("auth.ValidatePassword"),
		goerr.WithField("min_length", 6),
	)
}
