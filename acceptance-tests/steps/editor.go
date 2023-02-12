package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
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

func theManuscriptsToBeReviewedAreTheFollowing(ctx context.Context, table *godog.Table) (context.Context, error) {
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/to-review")
	var manuscriptsToReview api.ManuscriptsToReviewDto
	ctx, err := helpers.Call(ctx, url, http.MethodGet, nil, &manuscriptsToReview)
	if err != nil {
		return ctx, fmt.Errorf("unable to get manuscripts to review: %v", err)
	}
	return ctx, helpers.ShouldMatch(manuscriptsToReview.Manuscripts, table)
}

func EditorSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I review positively the manuscript for "([^"]*)"$`, reviewPositivelyManuscript)
	ctx.Step(`^the manuscripts to be reviewed are the following$`, theManuscriptsToBeReviewedAreTheFollowing)
}
