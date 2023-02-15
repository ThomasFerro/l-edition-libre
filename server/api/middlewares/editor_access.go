package middlewares

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

// TODO: On devrait pouvoir s'en passer grâce au scope de l'history
func OnlyAvailableForEditor(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userID := r.Context().Value(contexts.UserIDContextKey).(application.UserID)
		app := r.Context().Value(contexts.ApplicationContextKey).(application.Application)
		isAnEditor, err := app.UserIsAnEditor(userID)
		if err != nil {
			helpers.ManageError(w, err)
			return r
		}
		if !isAnEditor {
			helpers.ManageError(w, commands.ManuscriptNotFound{})
			return r
		}
		return next(w, r)
	}
}
