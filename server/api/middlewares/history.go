package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func InjectHistory() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			app := ApplicationFromRequest(r)

			userID := r.Context().Value(contexts.UserIDContextKey).(application.UserID)
			manuscriptsHistory, err := manuscriptsHistory(app, userID)
			if err != nil {
				helpers.ManageError(w, err)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptsHistoryContextKey, manuscriptsHistory))
			if manuscriptID, found := TryGetManuscriptID(r); found {
				r = r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptHistoryContextKey, application.ToEvents(manuscriptsHistory[manuscriptID])))
			}
			r = r.WithContext(context.WithValue(r.Context(), contexts.UserHistoryContextKey, []events.Event{}))
			next(w, r)
		}
	}
}

func manuscriptsHistory(app application.Application, userID application.UserID) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	isEditor, err := app.UserIsAnEditor(userID)
	if err != nil {
		return nil, err
	}
	if isEditor {
		return app.ManuscriptsHistory.ForAll()
	}
	return app.ManuscriptsHistory.ForAllOfUser(userID)
}
