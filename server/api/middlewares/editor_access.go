package middlewares

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

// TODO: On devrait pouvoir s'en passer grâce au scope de l'history
func OnlyAvailableForEditor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contexts.UserIDContextKey).(application.UserID)
		app := r.Context().Value(contexts.ApplicationContextKey).(application.Application)
		isAnEditor, err := app.UserIsAnEditor(userID)
		if err != nil {
			helpers.ManageError(w, err)
			return
		}
		if !isAnEditor {
			helpers.ManageError(w, commands.ManuscriptNotFound{})
			return
		}
		next(w, r)
	}
}
