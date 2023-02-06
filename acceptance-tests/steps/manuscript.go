package steps

import (
	testContext "acceptance-tests/test-context"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
)

type submitManuscriptRequest struct {
	ManuscriptName string `json:"manuscript_name"`
}
type submitManuscriptResponse struct {
	Id string `json:"id"`
}

func sumbitManuscript(ctx context.Context, manuscriptName string) (context.Context, error) {
	marshalled, err := json.Marshal(submitManuscriptRequest{
		ManuscriptName: manuscriptName,
	})
	if err != nil {
		return ctx, fmt.Errorf("unable to marshal the manuscript submission command: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/api/manuscripts", "application/json", bytes.NewReader(marshalled))
	if err != nil {
		return ctx, fmt.Errorf("unable to submit the manuscript - post error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ctx, fmt.Errorf("unable to submit the manuscript - wrong response code: %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ctx, fmt.Errorf("unable to submit the manuscript - body read error: %v", err)
	}
	var newManuscript submitManuscriptResponse
	err = json.Unmarshal(body, &newManuscript)
	if err != nil {
		return ctx, fmt.Errorf("unable to submit the manuscript - body unmarshal error: %v (body: %v)", err, string(body))
	}
	newManuscriptID := application.MustParseManuscriptID(newManuscript.Id)
	return testContext.SetManuscriptID(ctx, manuscriptName, newManuscriptID), nil
}

// TODO: déplacer dans le code de prod ?
type HttpErrorMessage struct {
	Error string `json:"error"`
}

func cancelManuscriptSubmission(ctx context.Context, manuscriptName string) (context.Context, error) {
	ctx, manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v/cancel-submission", manuscriptID.String())
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return ctx, fmt.Errorf("unable to cancel manuscript submission - post error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// TODO: Passer cette logique dans un handleError
		tags, ok := ctx.Value(testContext.TagsKey{}).([]*messages.PickleTag)
		if !ok {
			return ctx, fmt.Errorf("unable to get scenario context: %v", err)
		}
		isAnErrorHandlingScenario := false
		for _, nextTag := range tags {
			if nextTag.Name == "@Error" {
				isAnErrorHandlingScenario = true
				break
			}
		}

		if isAnErrorHandlingScenario {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ctx, fmt.Errorf("unable to cancel manuscript submission - body read error: %v", err)
			}

			var httpErrorMessage HttpErrorMessage
			err = json.Unmarshal(body, &httpErrorMessage)
			if err != nil {
				return ctx, fmt.Errorf("unable to cancel manuscript submission - body unmarshal error: %v (body: %v)", err, string(body))
			}
			return context.WithValue(ctx, testContext.ErrorKey{}, httpErrorMessage.Error), nil
		}
		return ctx, fmt.Errorf("unable to cancel manuscript submission - wrong response code: %v", resp.StatusCode)
	}
	return ctx, nil
}

func manuscriptStatusShouldBe(ctx context.Context, manuscriptName string, expectedStatus domain.Status) (context.Context, error) {
	ctx, manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v", manuscriptID.String())
	response, err := http.Get(url)

	if err != nil {
		return ctx, fmt.Errorf("unable to get manuscript's status: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return ctx, fmt.Errorf("unable to get manuscript's status - wrong response code: %v", response.StatusCode)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ctx, fmt.Errorf("unable to get manuscript's status - body read error: %v", err)
	}

	var manuscript api.ManuscriptDto
	err = json.Unmarshal(body, &manuscript)
	if err != nil {
		return ctx, fmt.Errorf("unable to get manuscript's status - body unmarshal error: %v", err)
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

func ManuscriptSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`I submit a PDF manuscript for "(.+?)"`, sumbitManuscript)
	ctx.Step(`I submitted a PDF manuscript for "(.+?)"`, sumbitManuscript)
	ctx.Step(`I cancel the submission of "(.+?)"`, cancelManuscriptSubmission)
	ctx.Step(`submission of "(.+?)" was canceled`, cancelManuscriptSubmission)
	ctx.Step(`"(.+?)" is pending review from the editor`, shouldBePendingReview)
	ctx.Step(`submission of "(.+?)" is canceled`, shouldBeCanceled)
}
