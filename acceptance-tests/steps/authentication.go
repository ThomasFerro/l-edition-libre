package steps

import (
	"context"

	"github.com/cucumber/godog"
)

func authentifyAsWriter(ctx context.Context) (context.Context, error) {
	// TODO: créer un nouveau writer et s'authentifier avec
	return ctx, nil
}

func AuthenticationSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`I am an authentified writer`, authentifyAsWriter)
}
