package middlewares

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func PersistNewEvents(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		r = next(w, r)

		newEvents := r.Context().Value(contexts.NewEventsContextKey{})
		if newEvents == nil {
			return r
		}

		userID := UserIdFromRequest(r)
		contextualizedUserEvents := []application.ContextualizedEvent{}
		contextualizedManuscriptEvents := []application.ContextualizedEvent{}
		contextualizedPublicationEvents := []application.ContextualizedEvent{}
		for _, nextEvent := range newEvents.([]events.Event) {
			if _, isUserEvent := nextEvent.(events.UserEvent); isUserEvent {
				contextualizedUserEvents = append(contextualizedUserEvents, application.ContextualizedEvent{
					OriginalEvent: nextEvent,
					Context: application.EventContext{
						UserID: userID,
					},
				})
				continue
			}
			if _, isManuscriptEvent := nextEvent.(events.ManuscriptEvent); isManuscriptEvent {
				contextualizedManuscriptEvents = append(contextualizedManuscriptEvents, application.ContextualizedEvent{
					OriginalEvent: nextEvent,
					Context: application.EventContext{
						UserID: userID,
					},
				})
				continue
			}
			if _, isPublicationEvent := nextEvent.(events.PublicationEvent); isPublicationEvent {
				contextualizedPublicationEvents = append(contextualizedPublicationEvents, application.ContextualizedEvent{
					OriginalEvent: nextEvent,
					Context: application.EventContext{
						UserID: userID,
					},
				})
			}
		}
		ctx := r.Context()
		if len(contextualizedUserEvents) != 0 {
			usersHistory := ctx.Value(contexts.UsersHistoryContextKey{}).(application.UsersHistory)
			err := usersHistory.Append(ctx, contextualizedUserEvents)
			if err != nil {
				helpers.ManageError(w, err)
				return r
			}
		}
		if len(contextualizedManuscriptEvents) != 0 {
			manuscriptsHistory := ctx.Value(contexts.ManuscriptsHistoryContextKey{}).(application.ManuscriptsHistory)
			err := manuscriptsHistory.Append(ctx, contextualizedManuscriptEvents)
			if err != nil {
				helpers.ManageError(w, err)
				return r
			}
		}
		if len(contextualizedPublicationEvents) != 0 {
			publicationsHistory := ctx.Value(contexts.PublicationsHistoryContextKey{}).(application.PublicationsHistory)
			err := publicationsHistory.Append(ctx, contextualizedPublicationEvents)
			if err != nil {
				helpers.ManageError(w, err)
				return r
			}
		}
		return r
	}
}
