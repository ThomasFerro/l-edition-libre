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
	ctx, err := helpers.Call(ctx, helpers.HttpRequest{
		Url:    url,
		Method: http.MethodPost,
	})
	if err != nil {
		return ctx, fmt.Errorf("unable to review manuscript: %v", err)
	}
	return ctx, nil
}

func tableToManuscriptsToReview(table *godog.Table) []api.ManuscriptToReviewDto {
	tableData := helpers.ExtractData(table)
	returned := []api.ManuscriptToReviewDto{}
	for _, row := range tableData {
		returned = append(returned, api.ManuscriptToReviewDto{
			Title:  row["Title"],
			Author: row["Author"],
		})
	}
	return returned
}

func manuscriptsToReviewAreExpected(actual, expected []api.ManuscriptToReviewDto) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("manuscripts to review mismatch\nexpected:\t%v\nactual:\t\t%v", expected, actual)
	}
	for index, nextExpected := range expected {
		nextActual := actual[index]
		if nextActual.Author != nextExpected.Author || nextActual.Title != nextExpected.Title {
			return fmt.Errorf("manuscripts to review mismatch\nexpected:\t%v\nactual:\t\t%v", expected, actual)
		}
	}
	return nil
}

func theManuscriptsToBeReviewedAreTheFollowing(ctx context.Context, table *godog.Table) (context.Context, error) {
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/to-review")
	var manuscriptsToReview api.ManuscriptsToReviewDto
	ctx, err := helpers.Call(ctx, helpers.HttpRequest{
		Url:         url,
		Method:      http.MethodGet,
		ResponseDto: &manuscriptsToReview,
	})
	if err != nil {
		return ctx, fmt.Errorf("unable to get manuscripts to review: %v", err)
	}

	expected := tableToManuscriptsToReview(table)
	return ctx, manuscriptsToReviewAreExpected(manuscriptsToReview.Manuscripts, expected)
}

func theManuscriptsForCanBeDownloadedForReview(ctx context.Context, manuscriptName string) (context.Context, error) {
	ctx, manuscriptID := helpers.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v/file", manuscriptID.String())
	ctx, err := helpers.Call(ctx, helpers.HttpRequest{
		Url:    url,
		Method: http.MethodGet,
	})
	if err != nil {
		return ctx, fmt.Errorf("unable to review manuscript: %v", err)
	}
	return ctx, nil
}

func EditorSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I review positively the manuscript for "([^"]*)"$`, reviewPositivelyManuscript)
	ctx.Step(`^the manuscripts to be reviewed are the following$`, theManuscriptsToBeReviewedAreTheFollowing)
	ctx.Step(`^the manuscripts for "([^"]*)" can be downloaded for review$`, theManuscriptsForCanBeDownloadedForReview)
}
