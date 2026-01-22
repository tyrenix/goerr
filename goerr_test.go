package goerr_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyrenix/goerr"
)

// notFound is a test error
var notFound = errors.New("not_found")

func TestNew(t *testing.T) {
	t.Run("creates error from string", func(t *testing.T) {
		err := goerr.New("test error")
		require.NotNil(t, err)
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("creates error from error", func(t *testing.T) {
		original := errors.New("original")
		err := goerr.New(original)
		require.NotNil(t, err)
		assert.True(t, errors.Is(err, original))
	})

	t.Run("with fields", func(t *testing.T) {
		err := goerr.New("test",
			goerr.Field("user_id", 123),
			goerr.Field("action", "login"))

		goErr := goerr.FromError(err)
		val, ok := goErr.GetField("user_id")
		assert.True(t, ok)
		assert.Equal(t, 123, val)
	})

	t.Run("with kind", func(t *testing.T) {
		err := goerr.New("test", goerr.Kind(notFound))
		goErr := goerr.FromError(err)
		assert.Equal(t, notFound, goErr.Kind())
	})

	t.Run("nil returns nil", func(t *testing.T) {
		err := goerr.New(nil)
		assert.Nil(t, err)
	})
}

func TestWrap(t *testing.T) {
	t.Run("wraps goerr preserving kind", func(t *testing.T) {
		original := goerr.New("not_found", goerr.Kind(notFound))
		wrapped := goerr.Wrap(original, "failed to get user", goerr.Field("user_id", 123))

		goErr := goerr.FromError(wrapped)
		assert.Equal(t, notFound, goErr.Kind())

		val, ok := goErr.GetField("user_id")
		assert.True(t, ok)
		assert.Equal(t, 123, val)
	})

	t.Run("doesn't mutate original", func(t *testing.T) {
		original := goerr.New("original")
		goerr.Wrap(original, "wrapped")

		assert.Equal(t, "original", original.Error())
	})

	t.Run("prepends context", func(t *testing.T) {
		err1 := goerr.New("base")
		err2 := goerr.Wrap(err1, "context1")
		err3 := goerr.Wrap(err2, "context2")

		details := goerr.FromError(err3).Details()
		assert.Contains(t, details, "context2")
		assert.Contains(t, details, "context1")
		assert.Contains(t, details, "base")
	})

	t.Run("wraps other errors", func(t *testing.T) {
		e0 := goerr.New(notFound, goerr.Field("user", 123), goerr.Kind(notFound))
		e1 := goerr.Wrap(e0, "failed to get user", goerr.Field("caller", "http"))
		e2 := goerr.Wrap(e1, "any step 1")
		e3 := goerr.Wrap(e2, "any step 2")

		assert.ErrorIs(t, e3, notFound)
		assert.ErrorIs(t, e3, e0)
		assert.ErrorIs(t, e3, e2)
		assert.Equal(t, notFound, goerr.FromError(e3).Kind())
	})
}

func TestDetails(t *testing.T) {
	t.Run("formats with fields", func(t *testing.T) {
		err := goerr.New("base_error",
			goerr.Field("user_id", 123),
			goerr.Field("action", "login"))

		details := goerr.FromError(err).Details()
		// check all fields exists
		assert.Contains(t, details, "user_id=123")
		assert.Contains(t, details, "action=login")
	})

	t.Run("formats wrapped chain", func(t *testing.T) {
		err1 := goerr.New("db error")
		err2 := goerr.Wrap(err1, "repo failed", goerr.Field("id", 1))
		err3 := goerr.Wrap(err2, "service failed")

		details := goerr.FromError(err3).Details()
		assert.Contains(t, details, "db error")
		assert.Contains(t, details, "repo failed")
		assert.Contains(t, details, "id=1")
	})
}

func TestGetField(t *testing.T) {
	t.Run("finds field in current error", func(t *testing.T) {
		err := goerr.New("test", goerr.Field("key", "value"))
		goErr := goerr.FromError(err)

		val, ok := goErr.GetField("key")
		assert.True(t, ok)
		assert.Equal(t, "value", val)
	})

	t.Run("finds field in wrapped error", func(t *testing.T) {
		err1 := goerr.New("base", goerr.Field("deep", "value"))
		err2 := goerr.Wrap(err1, "wrapper")
		goErr := goerr.FromError(err2)

		val, ok := goErr.GetField("deep")
		assert.True(t, ok)
		assert.Equal(t, "value", val)
	})

	t.Run("returns false for missing field", func(t *testing.T) {
		err := goerr.New("test")
		goErr := goerr.FromError(err)

		_, ok := goErr.GetField("missing")
		assert.False(t, ok)
	})
}
