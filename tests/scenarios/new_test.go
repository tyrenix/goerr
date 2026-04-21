package scenarios_test

import "github.com/daanila01/goerr/v3"

func (s *ErrorScenarioSuite) TestNew() {
	err := goerr.NewWithSpec("user not found", s.notFoundSpec.Code, s.notFoundSpec.Kind)

	s.Require().Error(err)

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)
	s.Equal("user not found", goErr.Error())
	s.Equal(s.notFoundSpec, goErr.Spec())
	s.Equal(goerr.Code("user.not_found"), goErr.Code())
	s.Equal(goerr.KindNotFound, goErr.Kind())
}

func (s *ErrorScenarioSuite) TestNew_WithoutSpec() {
	err := goerr.New("plain error")

	s.Require().Error(err)

	goErr, ok := goerr.AsError(err)
	s.Require().True(ok)
	s.Equal("plain error", goErr.Error())
	s.Equal(goerr.Spec{}, goErr.Spec())

	_, ok = goerr.CodeOf(err)
	s.False(ok)

	_, ok = goerr.KindOf(err)
	s.False(ok)
}
