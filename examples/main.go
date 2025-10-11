package main

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr"
)

func main() {
	one := errors.New("one")
	two := errors.New("two")
	err := goerr.New(one, two, "three", goerr.WithField("field", 100000), goerr.WithHTTPCode(500))

	fmt.Print("\n\n")
	fmt.Println("one is:", errors.Is(err, one))
	fmt.Println("two is:", errors.Is(err, two))
	fmt.Printf("short error: %s\n", err.Error())
	fmt.Println("full error:", err)

	wrapErr := goerr.New(err, "wrapped error")
	fmt.Print("\n\n")
	fmt.Println("wrapped short error:", wrapErr.Error())
	fmt.Println("wrapped full error:", wrapErr)
	fmt.Println("wrapped is:", errors.Is(wrapErr, err))
	fmt.Println("wrapped http code:", wrapErr.(*goerr.Error).HTTPCode())

	wrapGoErr := goerr.New("err1", goerr.New("err2", err))
	fmt.Print("\n\n")
	fmt.Println("wrapped go error:", wrapGoErr.Error())
	fmt.Println("wrapped full go error:", wrapGoErr)
	fmt.Println("wrapped http code:", wrapGoErr.(*goerr.Error).HTTPCode())
	fmt.Println("wrapped is:", errors.Is(wrapGoErr, err))
}
