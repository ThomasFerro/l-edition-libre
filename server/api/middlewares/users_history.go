package middlewares

import (
	"context"
	"errors"
	"net/http"

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
			r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedUserHistoryContextKey{}, func(ctx context.Context) ([]application.ContextualizedEvent, error) {
				return nil, errors.New("user not found")
			}))
			return next(w, r)
		}

		getUserHistory := func(ctx context.Context) ([]application.ContextualizedEvent, error) {
			return usersHistory.For(userID)
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedUserHistoryContextKey{}, getUserHistory))

		return next(w, r)
	}
}

func usersHistoryFromContext(ctx context.Context) application.UsersHistory {
	return ctx.Value(contexts.UsersHistoryContextKey{}).(application.UsersHistory)
}
