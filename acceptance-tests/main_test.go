package test

import (
	"acceptance-tests/steps"
	testContext "acceptance-tests/test-context"
	"fmt"
	"testing"

	"github.com/ThomasFerro/l-edition-libre/api"

	"github.com/go-bdd/gobdd"
)

func TestScenarios(t *testing.T) {
	go api.Start()
	suite := gobdd.NewSuite(t, gobdd.WithAfterScenario(func(ctx gobdd.Context) {
		unhandledError, err := ctx.Get(testContext.ErrorKey{})
		// TODO: Ne fonctionne pas? le test devrait planter car on n'enlève pas l'erreur du context
		fmt.Printf("\n\n%v\n\n\n", unhandledError)
		if err == nil && unhandledError != nil {
			t.Fatalf("unhandled error in tests: %v", unhandledError)
		}
	}))

	steps.AuthenticationSteps(suite)
	steps.ManuscriptSteps(suite)
	steps.ErrorSteps(suite)

	suite.Run()
}
