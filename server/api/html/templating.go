package html

import (
	"embed"
	"html/template"
	"net/http"
	"sort"

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

func WithFiles(files ...string) WithFilesOption {
	return WithFilesOption{
		Files: files,
	}
}

type WithFilesOption struct {
	Files []string
}

func (o WithFilesOption) Apply(t *template.Template) (*template.Template, error) {
	return t.ParseFS(templates, o.Files...)
}

type WithFuncs struct {
	Funcs template.FuncMap
}

func (o WithFuncs) Apply(t *template.Template) (*template.Template, error) {
	return t.Funcs(o.Funcs), nil
}

func RespondWithLayoutTemplate(w http.ResponseWriter, r *http.Request, data interface{}, options ...TemplateOption) *http.Request {
	files := WithFiles(
		"layout.gohtml",
		"authentication.gohtml",
	)
	return RespondWithTemplate(w, r, data, "layout", append(options, files)...)
}

func RespondWithErrorFragment(w http.ResponseWriter, r *http.Request, target string, err error) *http.Request {
	w.Header().Add("HX-Retarget", target)
	w.Header().Add("HX-Reswap", "beforeend")
	errorMessage := helpers.ExtractErrorMessage(err)
	return RespondWithTemplate(w, r, errorMessage.Error, "error-fragment", WithFiles("error-fragment.gohtml"))
}

func sortTemplateOption(options ...TemplateOption) []TemplateOption {
	sort.Slice(options, func(optionAIndex, _ int) bool {
		optionA := options[optionAIndex]
		switch optionA.(type) {
		case WithFilesOption:
			return false
		}
		return true
	})
	return options
}

func RespondWithTemplate(w http.ResponseWriter, r *http.Request, data interface{}, templateName string, options ...TemplateOption) *http.Request {
	t := template.New("")
	options = sortTemplateOption(options...)
	for _, option := range options {
		var err error
		t, err = option.Apply(t)
		if err != nil {
			slog.Error("template option error", err)
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
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
