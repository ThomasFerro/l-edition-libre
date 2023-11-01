package middlewares

import (
	"context"
	"net/http"

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
		userId, err := extractUserIDFromCookie(r)
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

func extractUserIDFromCookie(r *http.Request) (application.UserID, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}

	return auth0.ExtractUserID(cookie.Value)
}
