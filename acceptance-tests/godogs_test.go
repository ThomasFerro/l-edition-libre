package godogs_test

import (
	"acceptance-tests/helpers"
	"acceptance-tests/steps"
	"context"
	"fmt"
	"testing"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/cucumber/godog"
)

func InitializeTestSuite(*godog.TestSuiteContext) {
	go api.Start()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return context.WithValue(ctx, helpers.TagsKey{}, sc.Tags), nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		unhandledError, ok := ctx.Value(helpers.ErrorKey{}).(string)
		if ok {
			return ctx, fmt.Errorf("unhandled error in tests: %v", unhandledError)
		}
		return ctx, nil
	})

	steps.AuthenticationSteps(ctx)
	steps.ManuscriptSteps(ctx)
	steps.EditorSteps(ctx)
	steps.ErrorSteps(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: false,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
