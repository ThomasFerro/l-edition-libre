package steps

import (
	testContext "acceptance-tests/test-context"
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func errorThrown(ctx context.Context, errorType string) (context.Context, error) {
	scenarioError := ctx.Value(testContext.ErrorKey{})
	if scenarioError != errorType {
		return ctx, fmt.Errorf("expected a scenario error of type %v but got %v", errorType, scenarioError)
	}
	return context.WithValue(ctx, testContext.ErrorKey{}, nil), nil
}

func ErrorSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`the error "(.+?)" is thrown`, errorThrown)
}
