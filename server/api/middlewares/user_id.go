package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers/auth0"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"golang.org/x/exp/slog"
)

func TryGetUserIdFromRequest(r *http.Request) (contexts.UserID, bool) {
	value := r.Context().Value(contexts.UserIDContextKey{})
	if value == nil {
		return "", false
	}
	return value.(contexts.UserID), true
}

func UserIdFromRequest(r *http.Request) contexts.UserID {
	return r.Context().Value(contexts.UserIDContextKey{}).(contexts.UserID)
}

func TryExtractingUserID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userId, err := ExtractUserIDFromCookie(r)
		if err != nil || userId == "" {
			return next(w, r)
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey{}, userId))

		return next(w, r)
	}
}
func EnsureUserIsAuthenticatedAndExtractUserID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userId, err := ExtractUserIDFromCookie(r)
		if err != nil {
			slog.Warn("Unable to extract user id: %v", err)
			http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
			return r
		}
		if userId == "" {
			return next(w, r)
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey{}, userId))

		return next(w, r)
	}
}

func ExtractUserIDFromCookie(r *http.Request) (contexts.UserID, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}

	return auth0.ExtractUserID(cookie.Value)
}
