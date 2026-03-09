package scenarios_test

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v2"
)

func (s *ErrorScenarioSuite) TestHelpers() {
	err := goerr.Wrap(
		goerr.Wrap(
			errors.New("sql: no rows"),
			"execute query",
			goerr.WithSpec(s.notFoundSpec),
			goerr.WithField("user_id", 42),
			goerr.WithField("source", "repo"),
		),
		"get user",
		goerr.WithOp("userservice.Get"),
		goerr.WithField("source", "service"),
	)

	code, ok := goerr.CodeOf(err)
	s.Require().True(ok)
	s.Equal(goerr.Code("user.not_found"), code)

	kind, ok := goerr.KindOf(err)
	s.Require().True(ok)
	s.Equal(goerr.KindNotFound, kind)

	field, ok := goerr.FieldOf(err, "user_id")
	s.Require().True(ok)
	s.Equal(42, field)

	fields := goerr.AllFields(err)
	s.Require().NotNil(fields)
	s.Equal("service", fields["source"])
	s.Equal("userservice.Get", fields["op"])
	s.Equal(42, fields["user_id"])
}

func (s *ErrorScenarioSuite) TestHelpers_WithStdlibWrapper() {
	err := fmt.Errorf("transport layer: %w", goerr.Wrap(
		s.permissionError,
		"check access",
		goerr.WithSpec(goerr.Define(goerr.CodeForbidden, goerr.KindForbidden)),
	))

	code, ok := goerr.CodeOf(err)
	s.Require().True(ok)
	s.Equal(goerr.CodeForbidden, code)

	kind, ok := goerr.KindOf(err)
	s.Require().True(ok)
	s.Equal(goerr.KindForbidden, kind)

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)
	s.Equal("check access", goErr.Message())
}

func (s *ErrorScenarioSuite) TestHelpers_PlainError() {
	err := errors.New("plain")

	_, ok := goerr.CodeOf(err)
	s.False(ok)

	_, ok = goerr.KindOf(err)
	s.False(ok)

	_, ok = goerr.FieldOf(err, "op")
	s.False(ok)

	s.Nil(goerr.AllFields(err))
}

func (s *ErrorScenarioSuite) TestFields_Copy() {
	err := goerr.New(
		"validation failed",
		goerr.WithField("attempt", 1),
	)

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)

	fields := goErr.Fields()
	fields["attempt"] = 2
	fields["new"] = "value"

	val, exists := goErr.Field("attempt")
	s.Require().True(exists)
	s.Equal(1, val)

	_, exists = goErr.Field("new")
	s.False(exists)
}
