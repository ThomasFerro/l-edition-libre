package api

import (
	_ "embed"
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

//go:embed html/index.go.html
var index string

type TemplateManuscriptDto struct {
	Name string
}
type IndexParameters struct {
	Manuscripts   []TemplateManuscriptDto
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
		if isCurrentlyAuthenticated {
			http.Redirect(w, r, "/manuscripts", http.StatusFound)
			return r
		}

		t, err := template.New("Not authentified").Parse(index)
		if err != nil {
			slog.Error("not authentified template parsing error", err)
			helpers.ManageError(w, err)
			return r
		}
		err = t.Execute(w, nil)
		if err != nil {
			slog.Error("not authentified template execution error", err)
			helpers.ManageError(w, err)
			return r
		}
		return r
	}
}
