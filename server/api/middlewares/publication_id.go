package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func TryGetPublicationIdFromRequest(r *http.Request) (application.PublicationID, bool) {
	value := r.Context().Value(contexts.PublicationIDContextKey{})
	if value == nil {
		return application.PublicationID{}, false
	}
	return value.(application.PublicationID), true
}

func GetPublicationID(r *http.Request) application.PublicationID {
	return r.Context().Value(contexts.PublicationIDContextKey{}).(application.PublicationID)
}

func SetPublicationID(r *http.Request, publicationID application.PublicationID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contexts.PublicationIDContextKey{}, publicationID))
}

func ExtractPublicationID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		publicationID, err := application.ParsePublicationID(helpers.FromUrlParams(r.Context(), ":publicationID"))

		if err != nil {
			helpers.ManageErrorAsJson(w, err)
			return r
		}
		r = SetPublicationID(r, publicationID)
		return next(w, r)
	}
}
