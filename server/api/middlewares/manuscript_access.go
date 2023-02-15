package middlewares

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func UserShouldHaveAccessToManuscript(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userID := r.Context().Value(contexts.UserIDContextKey).(application.UserID)
		manuscriptID := GetManuscriptID(r)
		app := r.Context().Value(contexts.ApplicationContextKey).(application.Application)
		haveAccess, err := app.UserHaveAccessToManuscript(userID, manuscriptID)
		if err != nil {
			helpers.ManageError(w, err)
			return r
		}
		if !haveAccess {
			helpers.ManageError(w, commands.ManuscriptNotFound{})
			return r
		}
		return next(w, r)
	}
}
