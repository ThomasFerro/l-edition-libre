package helpers

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
)

type ErrorKey struct{}
type TagsKey struct{}

type AuthentifiedUserToken struct{}

func SetAuthentifiedUserToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AuthentifiedUserToken{}, token)
}

func GetAuthentifiedUserToken(ctx context.Context) string {
	return ctx.Value(AuthentifiedUserToken{}).(string)
}

type AuthentifiedUser struct{}

func SetAuthentifiedUserID(ctx context.Context, userID application.UserID) context.Context {
	return context.WithValue(ctx, AuthentifiedUser{}, userID)
}

func GetAuthentifiedUserID(ctx context.Context) application.UserID {
	return ctx.Value(AuthentifiedUser{}).(application.UserID)
}

type UserNameByIDKey struct{}
type UserNameByID map[application.UserID]string

func getOrCreateUserNameByIDFromContext(ctx context.Context) (context.Context, UserNameByID, error) {
	userNameByID, ok := ctx.Value(UserNameByIDKey{}).(UserNameByID)
	if ok {
		return ctx, userNameByID, nil
	}
	newMap := UserNameByID{}
	return context.WithValue(ctx, UserNameByIDKey{}, newMap), newMap, nil
}

func GetAuthentifiedUserName(ctx context.Context) (context.Context, string) {
	userID := GetAuthentifiedUserID(ctx)
	ctx, userNameByID, err := getOrCreateUserNameByIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, userNameByID[userID]
}

func SetUserName(ctx context.Context, userID application.UserID, userName string) context.Context {
	ctx, userNameByID, err := getOrCreateUserNameByIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	userNameByID[userID] = userName
	return context.WithValue(ctx, UserNameByIDKey{}, userNameByID)
}

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
