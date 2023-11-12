package helpers

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
)

type ErrorKey struct{}
type TagsKey struct{}

type AuthenticatedUserToken struct{}

func SetAuthenticatedUserToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AuthenticatedUserToken{}, token)
}

func GetAuthenticatedUserToken(ctx context.Context) string {
	return ctx.Value(AuthenticatedUserToken{}).(string)
}

type AuthenticatedUser struct{}

func SetAuthenticatedUserID(ctx context.Context, userID contexts.UserID) context.Context {
	return context.WithValue(ctx, AuthenticatedUser{}, userID)
}

func GetAuthenticatedUserID(ctx context.Context) contexts.UserID {
	value := ctx.Value(AuthenticatedUser{})
	if value == nil {
		return ""
	}
	return value.(contexts.UserID)
}

type UserNameByIDKey struct{}
type UserNameByID map[contexts.UserID]string

func getOrCreateUserNameByIDFromContext(ctx context.Context) (context.Context, UserNameByID, error) {
	userNameByID, ok := ctx.Value(UserNameByIDKey{}).(UserNameByID)
	if ok {
		return ctx, userNameByID, nil
	}
	newMap := UserNameByID{}
	return context.WithValue(ctx, UserNameByIDKey{}, newMap), newMap, nil
}

func GetAuthenticatedUserName(ctx context.Context) (context.Context, string) {
	userID := GetAuthenticatedUserID(ctx)
	ctx, userNameByID, err := getOrCreateUserNameByIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, userNameByID[userID]
}

func SetUserName(ctx context.Context, userID contexts.UserID, userName string) context.Context {
	ctx, userNameByID, err := getOrCreateUserNameByIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	userNameByID[userID] = userName
	return context.WithValue(ctx, UserNameByIDKey{}, userNameByID)
}

type TokenByUserIdKey struct{}
type TokenByUserID map[contexts.UserID]string

func getOrCreateTokenByIDFromContext(ctx context.Context) (context.Context, TokenByUserID, error) {
	tokenByUserID, ok := ctx.Value(TokenByUserIdKey{}).(TokenByUserID)
	if ok {
		return ctx, tokenByUserID, nil
	}
	newMap := TokenByUserID{}
	return context.WithValue(ctx, TokenByUserIdKey{}, newMap), newMap, nil
}

func GetUserToken(ctx context.Context) (context.Context, string) {
	userID := GetAuthenticatedUserID(ctx)
	ctx, tokenByUserID, err := getOrCreateTokenByIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, tokenByUserID[userID]
}

func SetToken(ctx context.Context, userID contexts.UserID, token string) context.Context {
	ctx, tokenByUserID, err := getOrCreateTokenByIDFromContext(ctx)
	if err != nil {
		panic(err)
	}

	tokenByUserID[userID] = token
	return context.WithValue(ctx, TokenByUserIdKey{}, tokenByUserID)
}

type ManuscriptIdByNameKey struct{}
type ManuscriptIdByName map[string]contexts.ManuscriptID

func getOrCreateManuscriptIdByNameFromContext(ctx context.Context) (context.Context, ManuscriptIdByName, error) {
	manuscriptIdByName, ok := ctx.Value(ManuscriptIdByNameKey{}).(ManuscriptIdByName)
	if ok {
		return ctx, manuscriptIdByName, nil
	}
	newMap := ManuscriptIdByName{}
	return context.WithValue(ctx, ManuscriptIdByNameKey{}, newMap), newMap, nil
}

func GetManuscriptID(ctx context.Context, manuscriptName string) (context.Context, contexts.ManuscriptID) {
	ctx, manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, manuscriptIdByName[manuscriptName]
}

func SetManuscriptID(ctx context.Context, manuscriptName string, manuscriptID contexts.ManuscriptID) context.Context {
	ctx, manuscriptIdByName, err := getOrCreateManuscriptIdByNameFromContext(ctx)
	if err != nil {
		panic(err)
	}

	manuscriptIdByName[manuscriptName] = manuscriptID
	return context.WithValue(ctx, ManuscriptIdByNameKey{}, manuscriptIdByName)
}
