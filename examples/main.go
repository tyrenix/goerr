package main

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr"
)

func main() {
	one := errors.New("one")
	two := errors.New("two")
	err := goerr.New(one, two, "three", goerr.WithHTTPCode(500))

	fmt.Println("one is:", errors.Is(err, one))
	fmt.Println("two is:", errors.Is(err, two))
	fmt.Printf("short error: %s\n", err.Error())
	fmt.Println("full error:", err)

	wrapErr := goerr.New(err, "wrapped error")
	fmt.Println("wrapped short error:", wrapErr.Error())
	fmt.Println("wrapped full error:", wrapErr)
	fmt.Println("wrapped is:", errors.Is(wrapErr, err))
}
