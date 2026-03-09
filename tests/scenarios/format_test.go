package scenarios_test

import (
	"errors"
	"fmt"

	"github.com/tyrenix/goerr/v2"
)

func (s *ErrorScenarioSuite) TestDetails() {
	err := goerr.Wrap(
		goerr.Wrap(
			errors.New("sql: no rows in result set"),
			"execute user query",
			goerr.WithSpec(s.notFoundSpec),
			goerr.WithOp("userrepo.GetUser"),
		),
		"get user",
		goerr.WithOp("userservice.GetUser"),
	)

	s.Equal(
		"user.not_found (kind=not_found): get user (op=userservice.GetUser): execute user query (op=userrepo.GetUser): sql: no rows in result set",
		goerr.DetailsOf(err),
	)
	s.Equal(err.Error(), fmt.Sprintf("%v", err))
	s.Equal(goerr.DetailsOf(err), fmt.Sprintf("%+v", err))
	s.Equal(err.Error(), fmt.Sprintf("%s", err))
	s.Equal(fmt.Sprintf("%q", err.Error()), fmt.Sprintf("%q", err))
}
