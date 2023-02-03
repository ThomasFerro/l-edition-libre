package steps

import (
	testContext "acceptance-tests/test-context"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"github.com/go-bdd/gobdd"
)

func sumbitManuscript(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)

	err := application.Send(commands.SubmitManuscript{
		ManuscriptName: manuscriptName,
	})
	if err != nil {
		t.Errorf("unable to submit the manuscript: %v", err)
	}
}

func cancelManuscriptSubmission(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)

	// TODO: Se baser sur un ID généré à la step d'avant
	err := application.Send(commands.CancelManuscriptSubmission{
		ManuscriptName: manuscriptName,
	})
	if err != nil {
		t.Errorf("unable to cancel the manuscript's submission: %v", err)
	}
}

func shouldBePendingReview(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	application := testContext.GetApp(ctx)

	// TODO: Se baser sur un ID généré à la step d'avant
	manuscriptStatus, err := application.Query(queries.ManuscriptStatus{
		ManuscriptName: manuscriptName,
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

	// TODO: Se baser sur un ID généré à la step d'avant
	manuscriptStatus, err := application.Query(queries.ManuscriptStatus{
		ManuscriptName: manuscriptName,
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
