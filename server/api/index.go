package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/html"
	"github.com/ThomasFerro/l-edition-libre/api/router"
)

func handleIndexFuncs(serveMux *http.ServeMux) {
	routes := []router.Route{
		{
			Path:    "/",
			Method:  "GET",
			Handler: handleIndex(),
		},
		{
			Path:    "/error",
			Method:  "GET",
			Handler: handleError(),
		},
		{
			Path:    "/htmx",
			Method:  "GET",
			Handler: handleHtmx(),
		},
	}
	router.HandleRoutes(serveMux, routes)
}

type IndexParameters struct {
	Authenticated bool
}

func handleIndex() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		isCurrentlyAuthenticated := isAuthenticated(r)

		return html.RespondWithLayoutTemplate(w, r, IndexParameters{
			Authenticated: isCurrentlyAuthenticated,
		}, html.WithFiles("index.gohtml"))
	}
}

func handleError() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		return html.RespondWithLayoutTemplate(w, r, nil, html.WithFiles("error-page.gohtml"))
	}
}

func handleHtmx() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		return html.RespondWithStaticFile(w, r, "htmx.min.js")
	}
}
