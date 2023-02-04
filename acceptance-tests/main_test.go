package test

import (
	"acceptance-tests/steps"
	"testing"

	"github.com/ThomasFerro/l-edition-libre/api"

	"github.com/go-bdd/gobdd"
)

func TestScenarios(t *testing.T) {
	go api.Start()
	suite := gobdd.NewSuite(t)

	steps.AuthenticationSteps(suite)
	steps.ManuscriptSteps(suite)

	suite.Run()
}
