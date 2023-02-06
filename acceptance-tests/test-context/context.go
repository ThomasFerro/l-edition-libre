package testContext

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
)

type ErrorKey struct{}
type TagsKey struct{}

type ManuscriptIdByNameKey struct{}
type ManuscriptIdByName map[string]application.ManuscriptID

func getOrCreateManuscriptIdByNameFromContext(ctx context.Context) (context.Context, ManuscriptIdByName, error) {
	manuscriptIdByName, ok := ctx.Value(ManuscriptIdByNameKey{}).(ManuscriptIdByName)
	if ok {
		return ctx, manuscriptIdByName, nil
	}
	newMap := ManuscriptIdByName{}
	return context.WithValue(ctx, ManuscriptIdByNameKey{}, newMap), newMap, nil
}

func GetManuscriptID(ctx context.Context, manuscriptName string) (context.Context, application.ManuscriptID) {
	ctx, manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, manuscriptIdByName[manuscriptName]
}

func SetManuscriptID(ctx context.Context, manuscriptName string, manuscriptID application.ManuscriptID) context.Context {
	ctx, manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	manuscriptIdByName[manuscriptName] = manuscriptID
	return context.WithValue(ctx, ManuscriptIdByNameKey{}, manuscriptIdByName)
}
