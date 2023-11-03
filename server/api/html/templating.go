package html

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"golang.org/x/exp/slog"
)

// TODO: Obligé de les embed ?
//
//go:embed *.gohtml
var templates embed.FS

type TemplateOption interface {
	Apply(t *template.Template) (*template.Template, error)
}

type WithFiles struct {
	Files []string
}

func (o WithFiles) Apply(t *template.Template) (*template.Template, error) {
	return t.ParseFS(templates, o.Files...)
}

type WithFuncs struct {
	Funcs template.FuncMap
}

func (o WithFuncs) Apply(t *template.Template) (*template.Template, error) {
	return t.Funcs(o.Funcs), nil
}

func RespondWithIndexTemplate(w http.ResponseWriter, r *http.Request, data interface{}, options ...TemplateOption) *http.Request {
	files := WithFiles{
		Files: []string{
			"layout.gohtml",
			"authentication.gohtml",
		}}
	return RespondWithTemplate(w, r, data, "layout", append(options, files)...)
}

func RespondWithErrorTemplate(w http.ResponseWriter, r *http.Request, target string, err error) *http.Request {
	w.Header().Add("HX-Retarget", target)
	w.Header().Add("HX-Reswap", "beforeend")
	errorMessage := helpers.ExtractErrorMessage(err)
	return RespondWithTemplate(w, r, errorMessage.Error, "error", WithFiles{Files: []string{"error.gohtml"}})
}

func RespondWithTemplate(w http.ResponseWriter, r *http.Request, data interface{}, templateName string, options ...TemplateOption) *http.Request {
	t := template.New("")
	for _, option := range options {
		var err error
		t, err = option.Apply(t)
		if err != nil {
			// TODO: Une gestion d'erreur critique comme celle-ci en json
			slog.Error("template option error", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}
	}

	err := t.ExecuteTemplate(w, templateName, data)
	if err != nil {
		slog.Error("template execution error", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}
	return r
}
