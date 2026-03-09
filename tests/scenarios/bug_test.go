package scenarios_test

import (
	"fmt"

	"github.com/tyrenix/goerr/v2"
)

func (s *ErrorScenarioSuite) TestBug_StdlibWrapperMustNotLeakDetailsIntoErrorText() {
	err := fmt.Errorf("transport layer: %w", goerr.Wrap(
		s.permissionError,
		"check access",
		goerr.WithSpec(goerr.Define(goerr.CodeForbidden, goerr.KindForbidden)),
	))

	s.Equal(
		"transport layer: check access: permission denied",
		err.Error(),
	)
	s.Equal(err.Error(), fmt.Sprintf("%v", err))
}

func (s *ErrorScenarioSuite) TestBug_DetailsOfMustKeepStdlibWrapperContext() {
	err := fmt.Errorf("transport layer: %w", goerr.Wrap(
		s.permissionError,
		"check access",
		goerr.WithSpec(goerr.Define(goerr.CodeForbidden, goerr.KindForbidden)),
	))

	s.Equal(
		"forbidden (kind=forbidden): transport layer: check access: permission denied",
		goerr.DetailsOf(err),
	)
	s.Equal(
		"transport layer: check access: permission denied",
		fmt.Sprintf("%+v", err),
	)
}
