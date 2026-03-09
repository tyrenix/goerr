package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v2"
)

var UserNotFound = goerr.Define("user.not_found", goerr.KindNotFound)

func main() {
	err := getUser()
	if err == nil {
		return
	}

	fmt.Println("error:", err.Error())
	fmt.Printf("details: %v\n", err)
	fmt.Println("fields:", goerr.AllFields(err))
}

func getUser() error {
	err := repoGetUser()
	if err != nil {
		return goerr.Wrap(
			err,
			"get user",
			goerr.WithOp("userservice.Get"),
			goerr.WithField("source", "service"),
		)
	}

	return nil
}

func repoGetUser() error {
	err := sql.ErrNoRows
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return goerr.Wrap(
			err,
			"execute user query",
			goerr.WithSpec(UserNotFound),
			goerr.WithOp("userrepo.Get"),
			goerr.WithField("user_id", 42),
		)
	}

	return nil
}
