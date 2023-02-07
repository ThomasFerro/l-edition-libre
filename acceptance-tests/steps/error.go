package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func errorThrown(ctx context.Context, errorType string) (context.Context, error) {
	scenarioError := ctx.Value(helpers.ErrorKey{})
	if scenarioError != errorType {
		return ctx, fmt.Errorf("expected a scenario error of type %v but got %v", errorType, scenarioError)
	}
	return context.WithValue(ctx, helpers.ErrorKey{}, nil), nil
}

func ErrorSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`the error "(.+?)" is thrown`, errorThrown)
}
