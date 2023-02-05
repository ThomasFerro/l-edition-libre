package steps

import (
	testContext "acceptance-tests/test-context"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/domain"
	msgs "github.com/cucumber/messages-go/v12"
	"github.com/go-bdd/gobdd"
)

type submitManuscriptRequest struct {
	ManuscriptName string `json:"manuscript_name"`
}
type submitManuscriptResponse struct {
	Id string `json:"id"`
}

func sumbitManuscript(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	marshalled, err := json.Marshal(submitManuscriptRequest{
		ManuscriptName: manuscriptName,
	})
	if err != nil {
		t.Fatalf("unable to marshal the manuscript submission command: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/api/manuscripts", "application/json", bytes.NewReader(marshalled))
	if err != nil {
		t.Fatalf("unable to submit the manuscript - post error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("unable to submit the manuscript - wrong response code: %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to submit the manuscript - body read error: %v", err)
	}
	var newManuscript submitManuscriptResponse
	err = json.Unmarshal(body, &newManuscript)
	if err != nil {
		t.Fatalf("unable to submit the manuscript - body unmarshal error: %v (body: %v)", err, string(body))
	}
	newManuscriptID := application.MustParseManuscriptID(newManuscript.Id)
	testContext.SetManuscriptID(ctx, manuscriptName, newManuscriptID)
}

// TODO: déplacer dans le code de prod ?
type HttpErrorMessage struct {
	Error string `json:"error"`
}

func cancelManuscriptSubmission(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v/cancel-submission", manuscriptID.String())
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		t.Fatalf("unable to cancel manuscript submission - post error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// TODO: Passer cette logique dans un handleError
		value, err := ctx.Get(gobdd.ScenarioKey{})
		if err != nil {
			t.Fatalf("unable to get scenario context: %v", err)
		}
		scenario := value.(*msgs.GherkinDocument_Feature_Scenario)
		isAnErrorHandlingScenario := false
		for _, nextTag := range scenario.Tags {
			if nextTag.Name == "@Error" {
				isAnErrorHandlingScenario = true
				break
			}
		}

		if isAnErrorHandlingScenario {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to cancel manuscript submission - body read error: %v", err)
			}

			var httpErrorMessage HttpErrorMessage
			err = json.Unmarshal(body, &httpErrorMessage)
			if err != nil {
				t.Fatalf("unable to cancel manuscript submission - body unmarshal error: %v (body: %v)", err, string(body))
			}
			ctx.Set(testContext.ErrorKey{}, httpErrorMessage.Error)
		} else {
			t.Fatalf("unable to cancel manuscript submission - wrong response code: %v", resp.StatusCode)
		}
	}
}

func manuscriptStatusShouldBe(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string, expectedStatus domain.Status) {
	manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)
	url := fmt.Sprintf("http://localhost:8080/api/manuscripts/%v", manuscriptID.String())
	response, err := http.Get(url)

	if err != nil {
		t.Fatalf("unable to get manuscript's status: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Fatalf("unable to get manuscript's status - wrong response code: %v", response.StatusCode)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to get manuscript's status - body read error: %v", err)
	}

	var manuscript api.ManuscriptDto
	err = json.Unmarshal(body, &manuscript)
	if err != nil {
		t.Fatalf("unable to get manuscript's status - body unmarshal error: %v", err)
	}

	if manuscript.Status != string(expectedStatus) {
		t.Fatalf("manuscript should be pending review instead of %v", manuscript.Status)
	}
}

func shouldBePendingReview(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	manuscriptStatusShouldBe(t, ctx, manuscriptName, domain.PendingReview)
}

func shouldBeCanceled(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	manuscriptStatusShouldBe(t, ctx, manuscriptName, domain.Canceled)
}

func ManuscriptSteps(suite *gobdd.Suite) {
	suite.AddStep(`I submit a PDF manuscript for "(.+?)"`, sumbitManuscript)
	suite.AddStep(`I submitted a PDF manuscript for "(.+?)"`, sumbitManuscript)
	suite.AddStep(`I cancel the submission of "(.+?)"`, cancelManuscriptSubmission)
	suite.AddStep(`submission of "(.+?)" was canceled`, cancelManuscriptSubmission)
	suite.AddStep(`"(.+?)" is pending review from the editor`, shouldBePendingReview)
	suite.AddStep(`submission of "(.+?)" is canceled`, shouldBeCanceled)
}
