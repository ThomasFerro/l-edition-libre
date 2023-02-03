package steps

import (
	testContext "acceptance-tests/test-context"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"github.com/go-bdd/gobdd"
)

// TODO: Passer par l'API et non directement les commandes

func sumbitManuscript(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)

	returnedEvents, err := application.Send(commands.SubmitManuscript{
		ManuscriptName: manuscriptName,
	})
	if err != nil {
		t.Errorf("unable to submit the manuscript: %v", err)
	}
	// TODO: Ne plus se baser sur les events retournées mais sur le retour de l'API HTTP
	for _, nextEvent := range returnedEvents {
		manuscriptSubmittedEvent, foundExpectedEvent := nextEvent.(events.ManuscriptSubmitted)
		if foundExpectedEvent {
			testContext.SetManuscriptID(ctx, manuscriptName, manuscriptSubmittedEvent.CreatedManuscriptID)
		}
	}
}

func cancelManuscriptSubmission(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)
	manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)

	_, err := application.Send(commands.CancelManuscriptSubmission{
		ManuscriptID: manuscriptID,
	})
	if err != nil {
		t.Errorf("unable to cancel the manuscript's submission: %v", err)
	}
}

func shouldBePendingReview(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)
	manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)

	manuscriptStatus, err := application.Query(queries.ManuscriptStatus{
		ManuscriptID: manuscriptID,
	})

	if err != nil {
		t.Errorf("unable to get manuscript's status: %v", err)
	}
	if manuscriptStatus != domain.PendingReview {
		t.Errorf("manuscript should be pending review instead of %v", manuscriptStatus)
	}
}

func shouldBeCanceled(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)
	manuscriptID := testContext.GetManuscriptID(ctx, manuscriptName)

	manuscriptStatus, err := application.Query(queries.ManuscriptStatus{
		ManuscriptID: manuscriptID,
	})

	if err != nil {
		t.Errorf("unable to get manuscript's status: %v", err)
	}
	if manuscriptStatus != domain.Canceled {
		t.Errorf("manuscript should be canceled review instead of %v", manuscriptStatus)
	}
}

func ManuscriptSteps(suite *gobdd.Suite) {
	suite.AddStep(`I submit a PDF manuscript for "(.+?)"`, sumbitManuscript)
	suite.AddStep(`I submitted a PDF manuscript for "(.+?)"`, sumbitManuscript)
	suite.AddStep(`I cancel the submission of "(.+?)"`, cancelManuscriptSubmission)
	suite.AddStep(`"(.+?)" is pending review from the editor`, shouldBePendingReview)
	suite.AddStep(`"(.+?)"'s submission is canceled`, shouldBeCanceled)
}
