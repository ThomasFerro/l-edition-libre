package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func TryGetManuscriptID(r *http.Request) (application.ManuscriptID, bool) {
	value := r.Context().Value(contexts.ManuscriptIDContextKey{})
	if value == nil {
		return application.ManuscriptID{}, false
	}
	return value.(application.ManuscriptID), true
}

func GetManuscriptID(r *http.Request) application.ManuscriptID {
	return r.Context().Value(contexts.ManuscriptIDContextKey{}).(application.ManuscriptID)
}

func SetManuscriptID(r *http.Request, manuscriptID application.ManuscriptID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptIDContextKey{}, manuscriptID))
}

func ExtractManuscriptID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		manuscriptID, err := application.ParseManuscriptID(helpers.FromUrlParams(r.Context(), ":manuscriptID"))

		if err != nil {
			helpers.ManageErrorAsJson(w, err)
			return r
		}
		r = SetManuscriptID(r, manuscriptID)
		return next(w, r)
	}
}
