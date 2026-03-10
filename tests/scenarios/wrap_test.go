package scenarios_test

import (
	"errors"
	"fmt"
)

func (s *ErrorScenarioSuite) TestStdlibWrappingOnly() {
	err := fmt.Errorf("get user by id: %w", fmt.Errorf("execute user query: %w: %w", s.errUserNotFound, s.sqlNoRows))

	s.True(errors.Is(err, s.errUserNotFound))
	s.True(errors.Is(err, s.sqlNoRows))
	s.Equal(
		"get user by id: execute user query: user not found: sql: no rows in result set",
		err.Error(),
	)
}
