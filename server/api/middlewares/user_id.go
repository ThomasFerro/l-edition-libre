package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

const UserIDHeader = "X-User-Id"

func UserIdFromRequest(r *http.Request) application.UserID {
	return r.Context().Value(contexts.UserIDContextKey).(application.UserID)
}

func SetUserId(r *http.Request, userID application.UserID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey, userID))
}

func ExtractUserID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userId := r.Header.Get(UserIDHeader)
		if userId == "" {
			return r
		}
		parsed, err := application.ParseUserID(userId)
		if err != nil {
			http.Error(w, "Malformed user id", http.StatusBadRequest)
			return r
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey, parsed))

		return next(w, r)
	}
}
