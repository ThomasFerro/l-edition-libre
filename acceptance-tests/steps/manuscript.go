package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/cucumber/godog"
)

func aWriterSubmittedAManuscript(ctx context.Context, manuscriptName string) (context.Context, error) {
	ctx, err := authentifyAsWriter(ctx)
	if err != nil {
		return nil, err
	}
	return sumbitManuscript(ctx, manuscriptName)
}

func theWriterSubmittedAManuscript(ctx context.Context, writerName string, manuscriptName string) (context.Context, error) {
	ctx, err := authentifyAsWriterWithName(ctx, writerName)
	if err != nil {
		return nil, err
	}
	return sumbitManuscript(ctx, manuscriptName)
}

func sumbitManuscript(ctx context.Context, manuscriptName string) (context.Context, error) {
	var newManuscript api.SubmitManuscriptResponseDto
	ctx, authentifiedUserName := helpers.GetAuthentifiedUserName(ctx)
	ctx, err := helpers.Call(ctx, "http://localhost:8080/api/manuscripts", http.MethodPost, api.SubmitManuscriptRequestDto{
		Title:  manuscriptName,
		Author: authentifiedUserName,
	}, &newManuscript)
	if err != nil {
		return ctx, fmt.Errorf("unable to submit the manuscript: %v", err)
	}

	newManuscriptID := application.MustParseManuscriptID(newManuscript.Id)
	return helpers.SetManuscriptID(ctx, manuscriptName, newManuscriptID), nil
}

func cancelManuscriptSubmission(ctx context.Context, manuscriptName string) (context.Context, error) {
	ctx, manuscriptID := helpers.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v/cancel-submission", manuscriptID.String())

	ctx, err := helpers.Call(ctx, url, http.MethodPost, nil, nil)
	if err != nil {
		return ctx, fmt.Errorf("unable to cancel manuscript submission: %v", err)
	}
	return ctx, nil
}

func getManuscriptStatus(ctx context.Context, manuscriptName string) (context.Context, api.ManuscriptDto, error) {
	ctx, manuscriptID := helpers.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v", manuscriptID.String())
	var manuscript api.ManuscriptDto
	ctx, err := helpers.Call(ctx, url, http.MethodGet, nil, &manuscript)
	if err != nil {
		return ctx, api.ManuscriptDto{}, fmt.Errorf("unable to get manuscript's status: %v", err)
	}
	return ctx, manuscript, nil
}

func manuscriptStatusShouldBe(ctx context.Context, manuscriptName string, expectedStatus domain.Status) (context.Context, error) {
	ctx, manuscript, err := getManuscriptStatus(ctx, manuscriptName)
	if err != nil {
		return ctx, fmt.Errorf("cannot check manuscript's status: %v", err)
	}

	if manuscript.Status != string(expectedStatus) {
		return ctx, fmt.Errorf("manuscript should be pending review instead of %v", manuscript.Status)
	}
	return ctx, nil
}

func shouldBePendingReview(ctx context.Context, manuscriptName string) (context.Context, error) {
	return manuscriptStatusShouldBe(ctx, manuscriptName, domain.PendingReview)
}

func shouldBeCanceled(ctx context.Context, manuscriptName string) (context.Context, error) {
	return manuscriptStatusShouldBe(ctx, manuscriptName, domain.Canceled)
}

func isEventuallyPublished(ctx context.Context, manuscriptName string) (context.Context, error) {
	// TODO: Gérer l'asynchrone en testant l'état de la publication créée
	return manuscriptStatusShouldBe(ctx, manuscriptName, domain.Reviewed)
}

func tryGetStatus(ctx context.Context, manuscriptName string) (context.Context, error) {
	ctx, _, err := getManuscriptStatus(ctx, manuscriptName)
	return ctx, err
}

func ManuscriptSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a writer submitted a manuscript for "(.+?)"$`, aWriterSubmittedAManuscript)
	ctx.Step(`^the writer "(.+?)" submitted a manuscript for "(.+?)"$`, theWriterSubmittedAManuscript)
	ctx.Step(`I submit a manuscript for "(.+?)"`, sumbitManuscript)
	ctx.Step(`I submitted a manuscript for "(.+?)"`, sumbitManuscript)
	ctx.Step(`I cancel the submission of "(.+?)"`, cancelManuscriptSubmission)
	ctx.Step(`submission of "(.+?)" was canceled`, cancelManuscriptSubmission)
	ctx.Step(`"(.+?)" is pending review from the editor`, shouldBePendingReview)
	ctx.Step(`^"([^"]*)" is eventually published$`, isEventuallyPublished)
	ctx.Step(`submission of "(.+?)" is canceled`, shouldBeCanceled)
	ctx.Step(`I try to get the submission status of "(.+?)"`, tryGetStatus)
}
