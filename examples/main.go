package main

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr"
)

func main() {
	// notFound is a test error
	errNotFound := errors.New("not_found")
	errInternal := errors.New("internal_error")

	// test 1: Basic error creation with fields and kind
	fmt.Println("=== Test 1: Basic Error Creation ===")
	err1 := goerr.New("user.not_found",
		goerr.Kind(errNotFound),
		goerr.Field("user_id", 123),
		goerr.Field("action", "login"),
	)

	fmt.Printf("Error: %s\n", err1.Error())
	fmt.Printf("Details: %v\n", err1)
	fmt.Printf("Kind: %v\n", goerr.FromError(err1).Kind())

	val, ok := goerr.FromError(err1).GetField("user_id")
	fmt.Printf("Field 'user_id': %v (exists: %v)\n", val, ok)

	// test 2: Wrapping errors preserves kind and adds context
	fmt.Println("\n=== Test 2: Wrapping Errors ===")
	err2 := goerr.Wrap(err1, "failed to get user from database")

	fmt.Printf("Wrapped Error: %s\n", err2.Error())
	fmt.Printf("Wrapped Details: %v\n", err2)
	fmt.Printf("Kind preserved: %v\n", goerr.FromError(err2).Kind())
	fmt.Printf("Original field accessible: %v\n", func() bool {
		_, ok := goerr.FromError(err2).GetField("user_id")
		return ok
	}())

	// test 3: Multiple wrapping layers
	fmt.Println("\n=== Test 3: Multiple Wrapping Layers ===")
	err3 := goerr.Wrap(err2,
		"service layer failed",
		goerr.Field("request_id", "abc-123"),
	)

	fmt.Printf("Multi-wrapped Error: %s\n", err3.Error())
	fmt.Printf("Multi-wrapped Details: %v\n", err3)

	// both fields should be accessible
	userID, _ := goerr.FromError(err3).GetField("user_id")
	reqID, _ := goerr.FromError(err3).GetField("request_id")
	fmt.Printf("user_id: %v, request_id: %v\n", userID, reqID)

	// test 4: errors.Is compatibility
	fmt.Println("\n=== Test 4: errors.Is Compatibility ===")
	baseErr := errors.New("database connection failed")
	wrappedErr := goerr.New(baseErr)
	doubleWrapped := goerr.Wrap(wrappedErr, "repository error")

	fmt.Printf("errors.Is works: %v\n", errors.Is(doubleWrapped, baseErr))

	// test 5: Real-world scenario
	fmt.Println("\n=== Test 5: Real-World Scenario ===")

	// repository layer
	repoErr := goerr.Wrap(goerr.New("internal", goerr.Kind(errInternal)), "postgres connection timeout")

	// service layer adds context
	serviceErr := goerr.Wrap(repoErr,
		"failed to fetch user",
		goerr.Field("user_id", 456),
	)

	// handler layer adds more context
	handlerErr := goerr.Wrap(serviceErr,
		"GET /api/users/:id failed",
		goerr.Field("endpoint", "/api/users/456"),
		goerr.Field("method", "GET"),
	)

	fmt.Printf("Handler sees: %s\n", handlerErr.Error())
	fmt.Printf("Full context: %v\n", handlerErr)
	fmt.Printf("Error kind: %v\n", goerr.FromError(handlerErr).Kind())

	// test 6: No mutation of original error
	fmt.Println("\n=== Test 6: No Mutation ===")
	original := goerr.New("original_error", goerr.Field("count", 1))
	wrapped1 := goerr.Wrap(original, "first wrap")
	wrapped2 := goerr.Wrap(original, "second wrap")

	fmt.Printf("Original: %s\n", original.Error())
	fmt.Printf("Wrapped 1: %s\n", wrapped1.Error())
	fmt.Printf("Wrapped 2: %s\n", wrapped2.Error())
	fmt.Printf("Original unchanged: %v\n", original.Error() == "original_error")
}
