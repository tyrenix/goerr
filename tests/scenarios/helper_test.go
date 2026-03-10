package scenarios_test

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v3"
)

func (s *ErrorScenarioSuite) TestHelpers_WithBusinessAndTechnicalErrors() {
	err := fmt.Errorf("execute user query: %w: %w", s.errUserNotFound, s.sqlNoRows)

	code, ok := goerr.CodeOf(err)
	s.Require().True(ok)
	s.Equal(goerr.Code("user.not_found"), code)
	s.True(goerr.CodeIs(err, s.notFoundSpec.Code))

	kind, ok := goerr.KindOf(err)
	s.Require().True(ok)
	s.Equal(goerr.KindNotFound, kind)
	s.True(goerr.KindIs(err, goerr.KindNotFound))

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)
	s.Equal("user not found", goErr.Error())

	s.True(errors.Is(err, s.errUserNotFound))
	s.True(errors.Is(err, s.sqlNoRows))
}

func (s *ErrorScenarioSuite) TestHelpers_NearestBusinessErrorWins() {
	repoErr := fmt.Errorf("execute user query: %w: %w", s.errUserNotFound, s.sqlNoRows)
	serviceErr := fmt.Errorf("authorize request: %w: %w", s.errUnauthorized, repoErr)

	code, ok := goerr.CodeOf(serviceErr)
	s.Require().True(ok)
	s.Equal(s.unauthorizedSpec.Code, code)

	kind, ok := goerr.KindOf(serviceErr)
	s.Require().True(ok)
	s.Equal(goerr.KindUnauthorized, kind)

	goErr, ok := goerr.AsError(serviceErr)
	s.Require().True(ok)
	s.Equal("unauthorized", goErr.Error())

	s.True(errors.Is(serviceErr, s.errUnauthorized))
	s.True(errors.Is(serviceErr, s.errUserNotFound))
	s.True(errors.Is(serviceErr, s.sqlNoRows))
}

func (s *ErrorScenarioSuite) TestHelpers_PlainError() {
	err := errors.New("plain")

	_, ok := goerr.CodeOf(err)
	s.False(ok)
	s.False(goerr.CodeIs(err, s.notFoundSpec.Code))

	_, ok = goerr.KindOf(err)
	s.False(ok)
	s.False(goerr.KindIs(err, goerr.KindNotFound))

	_, ok = goerr.AsError(err)
	s.False(ok)
}

func (s *ErrorScenarioSuite) TestErrorsIs_EquivalentSpec() {
	err := fmt.Errorf("execute user query: %w: %w", s.errUserNotFound, s.sqlNoRows)
	target := goerr.NewWithSpec("another message", s.notFoundSpec.Code, s.notFoundSpec.Kind)

	s.True(errors.Is(err, target))
}
