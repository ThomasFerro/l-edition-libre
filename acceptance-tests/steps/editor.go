package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

func reviewPositivelyManuscript(ctx context.Context, manuscriptName string) (context.Context, error) {
	ctx, manuscriptID := helpers.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v/review", manuscriptID.String())
	// ctx, err := helpers.Call(ctx, url, http.MethodPost, api.ReviewManuscriptRequestDto{}, nil)
	ctx, err := helpers.Call(ctx, url, http.MethodPost, nil, nil)
	if err != nil {
		return ctx, fmt.Errorf("unable to review manuscript: %v", err)
	}
	return ctx, nil
}

func EditorSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I review positively the manuscript for "([^"]*)"$`, reviewPositivelyManuscript)
}
