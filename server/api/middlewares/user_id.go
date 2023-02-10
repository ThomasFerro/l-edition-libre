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

func ExtractUserID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get(UserIDHeader)
		if userId == "" {
			return
		}
		parsed, err := application.ParseUserID(userId)
		if err != nil {
			http.Error(w, "Malformed user id", http.StatusBadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey, parsed))
		next(w, r)
	}
}
