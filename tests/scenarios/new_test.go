package scenarios_test

import "github.com/tyrenix/goerr/v2"

func (s *ErrorScenarioSuite) TestNew() {
	err := goerr.New(
		"password is too short",
		goerr.WithSpec(goerr.Define("password.too_short", goerr.KindInvalid)),
		goerr.WithOp("auth.ValidatePassword"),
		goerr.WithField("min_length", 6),
	)

	s.Require().Error(err)

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)
	s.Equal("password is too short", goErr.Message())
	s.Equal(goerr.Code("password.too_short"), goErr.Code())
	s.Equal(goerr.KindInvalid, goErr.Kind())

	val, exists := goErr.Field("min_length")
	s.Require().True(exists)
	s.Equal(6, val)
}
