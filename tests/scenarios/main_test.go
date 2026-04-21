package scenarios_test

import (
	"errors"
	"testing"

	"github.com/daanila01/goerr/v3"
	"github.com/stretchr/testify/suite"
)

type ErrorScenarioSuite struct {
	suite.Suite

	notFoundSpec     goerr.Spec
	unauthorizedSpec goerr.Spec
	invalidSpec      goerr.Spec
	errUserNotFound  error
	errUnauthorized  error
	sqlNoRows        error
}

func TestErrorScenarioSuite(t *testing.T) {
	suite.Run(t, new(ErrorScenarioSuite))
}

func (s *ErrorScenarioSuite) SetupTest() {
	s.notFoundSpec = goerr.Define("user.not_found", goerr.KindNotFound)
	s.unauthorizedSpec = goerr.Define("auth.unauthorized", goerr.KindUnauthorized)
	s.invalidSpec = goerr.Define("user.invalid", goerr.KindInvalid)
	s.errUserNotFound = goerr.NewWithSpec("user not found", s.notFoundSpec.Code, s.notFoundSpec.Kind)
	s.errUnauthorized = goerr.NewWithSpec("unauthorized", s.unauthorizedSpec.Code, s.unauthorizedSpec.Kind)
	s.sqlNoRows = errors.New("sql: no rows in result set")
}
