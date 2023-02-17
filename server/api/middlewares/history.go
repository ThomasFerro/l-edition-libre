package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func InjectUsersHistory(usersHistory application.UsersHistory) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.UsersHistoryContextKey, usersHistory))
			return next(w, r)
		}
	}
}

func InjectManuscriptsHistory(manuscriptsHistory application.ManuscriptsHistory) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptsHistoryContextKey, manuscriptsHistory))
			return next(w, r)
		}
	}
}

func InjectContextualizedManuscriptsHistory(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		manuscriptsHistory, err := manuscriptsHistory(r.Context())
		if err != nil {
			helpers.ManageError(w, err)
			return r
		}
		fmt.Printf("\n\n manuscriptsHistory %v \n\n", manuscriptsHistory)
		r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedManuscriptsHistoryContextKey, manuscriptsHistory))
		manuscriptID, found := TryGetManuscriptID(r)
		fmt.Printf("\n\n manuscriptID %v, found %v  ?? %v \n\n", manuscriptID, found, manuscriptsHistory[manuscriptID])
		if found {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedManuscriptHistoryContextKey, manuscriptsHistory[manuscriptID]))
		}
		return next(w, r)
	}
}

func InjectContextualizedUserHistory(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		usersHistory := usersHistoryFromContext(r.Context())
		userID, found := TryGetUserIdFromRequest(r)
		if !found {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedUserHistoryContextKey, []events.Event{}))
			return next(w, r)
		}
		userhistory, err := usersHistory.For(userID)
		if err != nil {
			helpers.ManageError(w, err)
			return r
		}
		r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedUserHistoryContextKey, application.ToEvents(userhistory)))
		return next(w, r)
	}
}

func manuscriptsHistoryFromContext(ctx context.Context) application.ManuscriptsHistory {
	return ctx.Value(contexts.ManuscriptsHistoryContextKey).(application.ManuscriptsHistory)
}

func usersHistoryFromContext(ctx context.Context) application.UsersHistory {
	return ctx.Value(contexts.UsersHistoryContextKey).(application.UsersHistory)
}

func manuscriptsHistory(ctx context.Context) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	isEditor, err := application.IsAnEditor(ctx)
	if err != nil {
		return nil, err
	}
	manuscriptsHistory := manuscriptsHistoryFromContext(ctx)
	userID := ctx.Value(contexts.UserIDContextKey).(application.UserID)
	if isEditor {
		return manuscriptsHistory.ForAll()
	}
	return manuscriptsHistory.ForAllOfUser(userID)
}
