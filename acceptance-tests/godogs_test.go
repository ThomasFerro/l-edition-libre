package godogs_test

import (
	"godogs/steps"
	"testing"

	"github.com/cucumber/godog"
)

func iEat(arg1 int) error {
	return godog.ErrPending
}

func thereAreGodogs(arg1 int) error {
	return godog.ErrPending
}

func thereShouldBeRemaining(arg1 int) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	steps.AuthenticationSteps(ctx)
	steps.ErrorSteps(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
