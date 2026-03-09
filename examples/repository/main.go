package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v2"
)

var (
	UserNotFound = goerr.Define("user.not_found", goerr.KindNotFound)
	Internal     = goerr.Define("internal", goerr.KindInternal)
)

func main() {
	err := getUser("42")
	if err == nil {
		return
	}

	fmt.Println("error:", err.Error())
	fmt.Printf("details: %+v\n", err)

	if code, ok := goerr.CodeOf(err); ok {
		fmt.Println("public code:", code)
	}

	if userID, ok := goerr.FieldOf(err, "user_id"); ok {
		fmt.Println("field user_id:", userID)
	}
}

func getUser(id string) error {
	err := executeQuery()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goerr.Wrap(
				err,
				"execute user query",
				goerr.WithSpec(UserNotFound),
				goerr.WithOp("userrepo.Get"),
				goerr.WithField("user_id", id),
			)
		}

		return goerr.Wrap(
			err,
			"execute user query",
			goerr.WithSpec(Internal),
			goerr.WithOp("userrepo.Get"),
			goerr.WithField("user_id", id),
		)
	}

	return nil
}

func executeQuery() error {
	return sql.ErrNoRows
}
