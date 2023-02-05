package steps

import "github.com/go-bdd/gobdd"

func authentifyAsWriter(t gobdd.StepTest, ctx gobdd.Context) {
	// TODO: créer un nouveau writer et s'authentifier avec
}

func AuthenticationSteps(suite *gobdd.Suite) {
	suite.AddStep(`I am an authentified writer`, authentifyAsWriter)
}
