package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func InjectPublicationsHistory(publicationsHistory application.PublicationsHistory) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.PublicationsHistoryContextKey{}, publicationsHistory))
			return next(w, r)
		}
	}
}

func InjectContextualizedPublicationHistory(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		publicationsHistory := publicationsHistoryFromContext(r.Context())
		publicationID, found := TryGetPublicationIdFromRequest(r)
		if !found {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedPublicationHistoryContextKey{}, []events.DecoratedEvent{}))
			return r
		}
		publicationhistory, err := publicationsHistory.For(publicationID)
		if err != nil {
			helpers.ManageError(w, err)
			return r
		}
		r = r.WithContext(context.WithValue(r.Context(), contexts.ContextualizedPublicationHistoryContextKey{}, mapHistory(publicationhistory)))
		return next(w, r)
	}
}

func publicationsHistoryFromContext(ctx context.Context) application.PublicationsHistory {
	return ctx.Value(contexts.PublicationsHistoryContextKey{}).(application.PublicationsHistory)
}
