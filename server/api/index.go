package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/html"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"golang.org/x/exp/slog"
)

func handleIndexFuncs(serveMux *http.ServeMux) {
	routes := []router.Route{
		{
			Path:    "/",
			Method:  "GET",
			Handler: handleIndex(),
		},
	}
	router.HandleRoutes(serveMux, routes)
}

type IndexParameters struct {
	Authenticated bool
}

func handleIndex() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		isCurrentlyAuthenticated, err := isAuthenticated(r)
		if err != nil {
			slog.Error("unable to check if currently authenticated", err)
			helpers.ManageError(w, err)
			return r
		}

		return html.RespondWithIndexTemplate(w, r, IndexParameters{
			Authenticated: isCurrentlyAuthenticated,
		}, "index.gohtml")
	}
}
