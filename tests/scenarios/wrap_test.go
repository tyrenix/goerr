package scenarios_test

import (
	"errors"

	"github.com/tyrenix/goerr/v2"
)

func (s *ErrorScenarioSuite) TestWrap() {
	cause := errors.New("db failed")

	err := goerr.Wrap(
		cause,
		"get user",
		goerr.WithOp("userservice.Get"),
	)

	s.Require().Error(err)
	s.Equal("get user: db failed", err.Error())
	s.ErrorIs(err, cause)
}

func (s *ErrorScenarioSuite) TestWrap_Nil() {
	err := goerr.Wrap(nil, "noop")
	s.NoError(err)
}

func (s *ErrorScenarioSuite) TestWrap_PreservesInnerSpecWhenOuterOnlyAddsContext() {
	err := goerr.Wrap(
		goerr.Wrap(
			errors.New("sql: no rows"),
			"execute query",
			goerr.WithSpec(s.notFoundSpec),
			goerr.WithOp("userrepo.Get"),
		),
		"get user",
		goerr.WithOp("userservice.Get"),
	)

	code, ok := goerr.CodeOf(err)
	s.Require().True(ok)
	s.Equal(goerr.Code("user.not_found"), code)

	kind, ok := goerr.KindOf(err)
	s.Require().True(ok)
	s.Equal(goerr.KindNotFound, kind)
}

func (s *ErrorScenarioSuite) TestWrap_OuterSpecOverridesInnerSpec() {
	err := goerr.Wrap(
		goerr.Wrap(
			errors.New("sql: no rows"),
			"execute query",
			goerr.WithSpec(s.notFoundSpec),
		),
		"map repo error",
		goerr.WithSpec(s.invalidSpec),
	)

	code, ok := goerr.CodeOf(err)
	s.Require().True(ok)
	s.Equal(goerr.Code("user.invalid"), code)

	kind, ok := goerr.KindOf(err)
	s.Require().True(ok)
	s.Equal(goerr.KindInvalid, kind)
}
