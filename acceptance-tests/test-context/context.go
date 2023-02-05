package testContext

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/go-bdd/gobdd/context"
)

type ErrorKey struct{}

type ManuscriptIdByNameKey struct{}
type ManuscriptIdByName map[string]application.ManuscriptID

func getOrCreateManuscriptIdByNameFromContext(ctx context.Context) (ManuscriptIdByName, error) {
	manuscriptIdByName, err := ctx.Get(ManuscriptIdByNameKey{})
	if err == nil {
		return manuscriptIdByName.(ManuscriptIdByName), nil
	}
	newMap := ManuscriptIdByName{}
	ctx.Set(ManuscriptIdByNameKey{}, newMap)
	return newMap, nil
}

func GetManuscriptID(ctx context.Context, manuscriptName string) application.ManuscriptID {
	manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return manuscriptIdByName[manuscriptName]
}

func SetManuscriptID(ctx context.Context, manuscriptName string, manuscriptID application.ManuscriptID) {
	manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	manuscriptIdByName[manuscriptName] = manuscriptID
	ctx.Set(ManuscriptIdByNameKey{}, manuscriptIdByName)
}
