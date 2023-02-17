package middlewares

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
)

func UserShouldHaveAccessToManuscript(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		haveAccess, err := application.UserHaveAccessToManuscript(r.Context())
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
