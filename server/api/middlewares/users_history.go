package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func InjectUsersHistory(usersHistory application.UsersHistory) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.UsersHistoryContextKey{}, usersHistory))
			return next(w, r)
		}
	}
}

func InjectContextualizedUserHistory(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		usersHistory := usersHistoryFromContext(r.Context())
		userID, found := TryGetUserIdFromRequest(r)
		if !found {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedUserHistoryContextKey{}, []application.ContextualizedEvent{}))
			return next(w, r)
		}
		userhistory, err := usersHistory.For(userID)
		if err != nil {
			helpers.ManageError(w, err)
			return r
		}
		r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedUserHistoryContextKey{}, userhistory))
		return next(w, r)
	}
}

func usersHistoryFromContext(ctx context.Context) application.UsersHistory {
	return ctx.Value(contexts.UsersHistoryContextKey{}).(application.UsersHistory)
}
