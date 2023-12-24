package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/html"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
)

func handleIndexFuncs(serveMux *http.ServeMux, usersHistory application.UsersHistory) {
	routes := []router.Route{
		{
			Path:    "/",
			Method:  "GET",
			Handler: handleIndex(),
			Middlewares: []middlewares.Middleware{
				middlewares.InjectContextualizedUserHistory,
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.TryExtractingUserID,
			},
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
	Editor        bool
}

func handleIndex() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		_, isCurrentlyAuthenticated := middlewares.TryGetUserIdFromRequest(r)
		isAnEditor, _ := application.IsAnEditor(r.Context())

		return html.RespondWithLayoutTemplate(w, r, IndexParameters{
			Authenticated: isCurrentlyAuthenticated,
			Editor:        isAnEditor,
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
