package api

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
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

//go:embed html/*
var templates embed.FS

type IndexParameters struct {
	Authenticated bool
}

func handleIndex() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		t, err := template.New("index").ParseFS(templates, "html/*.gohtml")
		if err != nil {
			slog.Error("index template parsing error", err)
			helpers.ManageError(w, err)
			return r
		}

		isCurrentlyAuthenticated, err := isAuthenticated(r)
		if err != nil {
			slog.Error("unable to check if currently authenticated", err)
			helpers.ManageError(w, err)
			return r
		}

		err = t.ExecuteTemplate(w, "index", IndexParameters{
			Authenticated: isCurrentlyAuthenticated,
		})
		if err != nil {
			slog.Error("index template execution error", err)
			helpers.ManageError(w, err)
			return r
		}
		return r
	}
}
