package testContext

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/go-bdd/gobdd/context"
)

type AppKey struct{}

func GetApp(ctx context.Context) *application.Application {
	app, err := ctx.Get(AppKey{})
	if err != nil {
		panic(err)
	}

	return app.(*application.Application)
}

type ManuscriptIdByNameKey struct{}
type ManuscriptIdByName map[string]events.ManuscriptID

func getOrCreateManuscriptIdByNameFromContext(ctx context.Context) (ManuscriptIdByName, error) {
	manuscriptIdByName, err := ctx.Get(ManuscriptIdByNameKey{})
	if err == nil {
		return manuscriptIdByName.(ManuscriptIdByName), nil
	}
	newMap := ManuscriptIdByName{}
	ctx.Set(ManuscriptIdByNameKey{}, newMap)
	return newMap, nil
}

func GetManuscriptID(ctx context.Context, manuscriptName string) events.ManuscriptID {
	manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return manuscriptIdByName[manuscriptName]
}

func SetManuscriptID(ctx context.Context, manuscriptName string, manuscriptID events.ManuscriptID) {
	manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	manuscriptIdByName[manuscriptName] = manuscriptID
	ctx.Set(ManuscriptIdByNameKey{}, manuscriptIdByName)
}
