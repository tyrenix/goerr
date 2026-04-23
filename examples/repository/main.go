package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/daanila01/goerr"
)

var (
	ErrUserNotFound = goerr.NewWithSpec("user not found", "user.not_found", goerr.KindNotFound)
	ErrInternal     = goerr.NewWithSpec("internal", "internal", goerr.KindInternal)
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
