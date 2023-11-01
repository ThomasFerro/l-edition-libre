package html

import (
	"embed"
	"html/template"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
)

// TODO: Obligé de les embed ?
//
//go:embed *.gohtml
var templates embed.FS

func RespondWithTemplate(w http.ResponseWriter, r *http.Request, data interface{}, specificFiles ...string) *http.Request {
	files := []string{
		"layout.gohtml",
		"authentication.gohtml",
	}
	for _, specificFile := range specificFiles {
		files = append(files, specificFile)
	}

	t, err := template.New("index").ParseFS(templates, files...)
	if err != nil {
		slog.Error("template parsing error", err)
		helpers.ManageError(w, err)
		return r
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		slog.Error("template execution error", err)
		helpers.ManageError(w, err)
		return r
	}
	return r
}
