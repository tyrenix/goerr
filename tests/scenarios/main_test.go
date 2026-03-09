package scenarios_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tyrenix/goerr/v2"
)

type ErrorScenarioSuite struct {
	suite.Suite

	notFoundSpec    goerr.Spec
	invalidSpec     goerr.Spec
	permissionError error
}

func TestErrorScenarioSuite(t *testing.T) {
	suite.Run(t, new(ErrorScenarioSuite))
}

func (s *ErrorScenarioSuite) SetupTest() {
	s.notFoundSpec = goerr.Define("user.not_found", goerr.KindNotFound)
	s.invalidSpec = goerr.Define("user.invalid", goerr.KindInvalid)
	s.permissionError = errors.New("permission denied")
}
