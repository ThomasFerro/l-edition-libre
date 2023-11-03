package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/utils"
)

func InjectManuscriptsHistory(manuscriptsHistory application.ManuscriptsHistory) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptsHistoryContextKey{}, manuscriptsHistory))
			return next(w, r)
		}
	}
}

func InjectContextualizedManuscriptsHistory(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		manuscriptsHistory, err := manuscriptsHistory(r.Context())
		if err != nil {
			helpers.ManageErrorAsJson(w, err)
			return r
		}
		r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedManuscriptsHistoryContextKey{}, mapHistories(manuscriptsHistory)))
		manuscriptID, found := TryGetManuscriptID(r)
		if found {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedManuscriptHistoryContextKey{}, mapHistory(manuscriptsHistory.Of(manuscriptID))))
		}
		return next(w, r)
	}
}

func manuscriptsHistoryFromContext(ctx context.Context) application.ManuscriptsHistory {
	return ctx.Value(contexts.ManuscriptsHistoryContextKey{}).(application.ManuscriptsHistory)
}

func manuscriptsHistory(ctx context.Context) (utils.OrderedMap[application.ManuscriptID, []application.ContextualizedEvent], error) {
	isEditor, err := application.IsAnEditor(ctx)
	if err != nil {
		return utils.OrderedMap[application.ManuscriptID, []application.ContextualizedEvent]{}, err
	}
	manuscriptsHistory := manuscriptsHistoryFromContext(ctx)
	if isEditor {
		return manuscriptsHistory.ForAll()
	}
	userID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
	return manuscriptsHistory.ForAllOfUser(userID)
}
