package middlewares

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/domain"
)

func UserShouldHaveAccessToManuscript(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		haveAccess, err := application.UserHaveAccessToManuscript(r.Context())
		if err != nil {
			helpers.ManageErrorAsJson(w, err)
			return r
		}
		if !haveAccess {
			helpers.ManageErrorAsJson(w, domain.ManuscriptNotFound{})
			return r
		}
		return next(w, r)
	}
}
