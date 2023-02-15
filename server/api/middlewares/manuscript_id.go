package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func TryGetManuscriptID(r *http.Request) (application.ManuscriptID, bool) {
	value := r.Context().Value(contexts.ManuscriptIDContextKey)
	if value == nil {
		return application.ManuscriptID{}, false
	}
	return value.(application.ManuscriptID), true
}

func GetManuscriptID(r *http.Request) application.ManuscriptID {
	return r.Context().Value(contexts.ManuscriptIDContextKey).(application.ManuscriptID)
}

func SetManuscriptID(r *http.Request, manuscriptID application.ManuscriptID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptIDContextKey, manuscriptID))
}

func ExtractManuscriptID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		manuscriptID, err := application.ParseManuscriptID(helpers.FromUrlParams(r.Context(), ":manuscriptID"))

		if err != nil {
			helpers.ManageError(w, err)
			return
		}
		r = SetManuscriptID(r, manuscriptID)
		next(w, r)
	}
}
