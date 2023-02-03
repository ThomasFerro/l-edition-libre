package test

import (
	"acceptance-tests/steps"
	testContext "acceptance-tests/test-context"
	"testing"

	"github.com/ThomasFerro/l-edition-libre/application"

	"github.com/go-bdd/gobdd"
	"github.com/go-bdd/gobdd/context"
)

func TestScenarios(t *testing.T) {
	suite := gobdd.NewSuite(t, gobdd.WithBeforeScenario(func(ctx context.Context) {
		app := application.NewApplication()
		ctx.Set(testContext.AppKey{}, &app)
	}))

	steps.AuthenticationSteps(suite)
	steps.ManuscriptSteps(suite)

	suite.Run()
}
