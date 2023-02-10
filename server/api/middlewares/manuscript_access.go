package middlewares

import (
	"fmt"
	"net/http"

	apiContext "github.com/ThomasFerro/l-edition-libre/api/context"
	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
)

func UserShouldHaveAccessToManuscript(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(apiContext.UserIDContextKey).(application.UserID)
		manuscriptID := r.Context().Value(apiContext.ManuscriptIDContextKey).(application.ManuscriptID)
		app := r.Context().Value(apiContext.ApplicationContextKey).(application.Application)
		haveAccess, err := app.UserHaveAccessToManuscript(userID, manuscriptID)
		if err != nil {
			helpers.ManageError(w, err)
			return
		}
		if !haveAccess {
			helpers.ManageError(w, commands.ManuscriptNotFound{})
			return
		}
		fmt.Printf("\n\n\nHEHE\n\n\n")
		next(w, r)
	}
}
