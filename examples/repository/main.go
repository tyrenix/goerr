package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v3"
)

var (
	ErrUserNotFound = goerr.New(
		"user not found",
		goerr.WithSpec(goerr.Define("user.not_found", goerr.KindNotFound)),
	)
	ErrInternal = goerr.New(
		"internal",
		goerr.WithSpec(goerr.Define("internal", goerr.KindInternal)),
	)
)

func main() {
	err := getUser("42")
	if err == nil {
		return
	}

	code, _ := goerr.CodeOf(err)
	kind, _ := goerr.KindOf(err)

	fmt.Println("error:", err)
	fmt.Println("code:", code)
	fmt.Println("kind:", kind)
	fmt.Println("is user not found:", errors.Is(err, ErrUserNotFound))
	fmt.Println("is sql no rows:", errors.Is(err, sql.ErrNoRows))
}

func getUser(id string) error {
	err := executeQuery(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("execute user query: %w: %w", ErrUserNotFound, err)
		}

		return fmt.Errorf("execute user query: %w: %w", ErrInternal, err)
	}

	return nil
}

func executeQuery(id string) error {
	_ = id
	return sql.ErrNoRows
}
