package scenarios_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tyrenix/goerr/v3"
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
	s.errUserNotFound = goerr.New("user not found", goerr.WithSpec(s.notFoundSpec))
	s.errUnauthorized = goerr.New("unauthorized", goerr.WithSpec(s.unauthorizedSpec))
	s.sqlNoRows = errors.New("sql: no rows in result set")
}
