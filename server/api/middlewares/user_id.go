package middlewares

import (
	"context"
	"net/http"

	apiContext "github.com/ThomasFerro/l-edition-libre/api/context"
	"github.com/ThomasFerro/l-edition-libre/application"
)

const UserIDHeader = "X-User-Id"

func UserIdFromRequest(r *http.Request) application.UserID {
	return r.Context().Value(apiContext.UserIDContextKey).(application.UserID)
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
		r = r.WithContext(context.WithValue(r.Context(), apiContext.UserIDContextKey, parsed))
		next(w, r)
	}
}
