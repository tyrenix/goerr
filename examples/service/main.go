package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v3"
)

var ErrUserNotFound = goerr.New(
	"user not found",
	goerr.WithSpec(goerr.Define("user.not_found", goerr.KindNotFound)),
)

func main() {
	err := getUserByID("42")
	if err == nil {
		return
	}

	code, _ := goerr.CodeOf(err)
	kind, _ := goerr.KindOf(err)

	fmt.Println("error:", err)
	fmt.Println("code:", code)
	fmt.Println("kind:", kind)
	fmt.Println("is not found:", errors.Is(err, ErrUserNotFound))
	fmt.Println("is sql no rows:", errors.Is(err, sql.ErrNoRows))
}

func getUserByID(id string) error {
	if err := repoGetUserByID(id); err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	return nil
}

func repoGetUserByID(id string) error {
	_ = id
	return fmt.Errorf("execute user query: %w: %w", ErrUserNotFound, sql.ErrNoRows)
}
