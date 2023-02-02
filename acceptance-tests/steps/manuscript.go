package steps

import (
	testContext "acceptance-tests/test-context"

	"github.com/ThomasFerro/l-edition-libre/commands"
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

func shouldBePendingReview(t gobdd.StepTest, ctx gobdd.Context, manuscriptName string) {
	// TODO
}

func ManuscriptSteps(suite *gobdd.Suite) {
	suite.AddStep(`I submit a PDF manuscript for "(.+?)"`, sumbitManuscript)
	suite.AddStep(`"(.+?)" is pending review from the editor`, shouldBePendingReview)
}
