package scenarios_test

import (
	"fmt"

	"github.com/daanila01/goerr/v3"
)

func (s *ErrorScenarioSuite) TestAsError_WithNestedFmtErrorf() {
	err := fmt.Errorf("service layer: %w", fmt.Errorf("repo layer: %w: %w", s.errUserNotFound, s.sqlNoRows))

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)
	s.Equal("user not found", goErr.Error())

	code, ok := goerr.CodeOf(err)
	s.Require().True(ok)
	s.Equal(s.notFoundSpec.Code, code)
}
