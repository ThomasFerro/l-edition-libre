package steps

import (
	testContext "acceptance-tests/test-context"

	"github.com/go-bdd/gobdd"
)

func errorThrown(t gobdd.StepTest, ctx gobdd.Context, errorType string) {
	scenarioError, err := ctx.Get(testContext.ErrorKey{})
	if err != nil {
		t.Fatalf("cannot get error from test context: %v", err)
	}
	if scenarioError != errorType {
		t.Fatalf("expected a scenario error of type %v but got %v", errorType, scenarioError)
	}
	// TODO: Sortir l'erreur du context
}

func ErrorSteps(suite *gobdd.Suite) {
	suite.AddStep(`the error "(.+?)" is thrown`, errorThrown)
}
