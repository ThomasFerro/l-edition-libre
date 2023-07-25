package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/api/helpers/auth0"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"golang.org/x/exp/slog"
)

func TryGetUserIdFromRequest(r *http.Request) (application.UserID, bool) {
	value := r.Context().Value(contexts.UserIDContextKey{})
	if value == nil {
		return "", false
	}
	return value.(application.UserID), true
}

func UserIdFromRequest(r *http.Request) application.UserID {
	return r.Context().Value(contexts.UserIDContextKey{}).(application.UserID)
}

func ExtractUserID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userId, err := extractUserIDFromJwt(r)
		if err != nil {
			slog.Warn("Unable to extract user id: %v", err)
			http.Error(w, "Unable to extract user id", http.StatusBadRequest)
			return r
		}
		if userId == "" {
			return next(w, r)
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey{}, userId))

		return next(w, r)
	}
}

func extractUserIDFromJwt(r *http.Request) (application.UserID, error) {
	token, err := extractBearerToken(r.Header)
	if err != nil {
		return "", err
	}

	return auth0.ExtractUserID(token)
}

func extractBearerToken(header http.Header) (string, error) {
	authorizationHeader := header.Get("Authorization")
	split := strings.Split(authorizationHeader, "Bearer ")
	if len(split) != 2 {
		return "", errors.New("Unable to extract bearer")
	}
	return split[1], nil
}
