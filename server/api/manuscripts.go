package api

import (
	"io"
	"net/http"
	"net/url"
	"text/template"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/html"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/ports"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type WriterManuscriptDto struct {
	Title  string                  `json:"title"`
	Status domain.ManuscriptStatus `json:"status"`
}
type WriterManuscriptsDto struct {
	Authenticated bool
	Manuscripts   []WriterManuscriptDto
}

var translateStatus = func(status domain.ManuscriptStatus) string {
	switch status {
	case domain.PendingReview:
		return "en attente de revue"
	case domain.Reviewed:
		return "validé"
	case domain.Canceled:
		return "soumission annulée"
	default:
		return "status introuvable"
	}
}
var transateStatusFuncMap = template.FuncMap{
	"translateStatus": translateStatus,
}

func handleGetManuscripts(w http.ResponseWriter, r *http.Request) *http.Request {
	slog.Info("writer manuscripts request")
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.Query(r.Context(), queries.WriterManuscripts{})
	if err != nil {
		slog.Error("writer manuscripts query error", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}
	manuscripts, castedSuccessfuly := queryResult.([]domain.Manuscript)
	if !castedSuccessfuly {
		slog.Error("writer manuscripts query result casting error", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}

	dto := WriterManuscriptsDto{
		Manuscripts:   []WriterManuscriptDto{},
		Authenticated: true,
	}
	for _, manuscript := range manuscripts {
		dto.Manuscripts = append(dto.Manuscripts, WriterManuscriptDto{
			Title:  manuscript.Title,
			Status: manuscript.Status,
		})
	}

	return html.RespondWithLayoutTemplate(w, r, dto, html.WithFuncs{
		Funcs: transateStatusFuncMap,
	}, html.WithFiles("manuscripts.gohtml", "manuscript-item.gohtml"))
}

type SubmitManuscriptRequestDto struct {
	Title    string
	Author   string
	File     io.Reader
	FileName string
}

type ManuscriptCreationError struct{}

func (e ManuscriptCreationError) Error() string {
	return "Une erreur est survenue à la soumission de votre manuscrit"
}

func handleManuscriptCreation(w http.ResponseWriter, r *http.Request) *http.Request {
	file, _, err := r.FormFile("file")
	if err != nil {
		slog.Error("manuscript creation request file reading error", err)
		return html.RespondWithErrorFragment(w, r, "#new-manuscript-form", ManuscriptCreationError{})
	}
	dto := SubmitManuscriptRequestDto{
		Title:  r.FormValue("title"),
		Author: r.FormValue("author"),
		File:   file,
	}

	slog.Info("receiving manuscript creation request", "body", dto)

	newManuscriptID := application.NewManuscriptID()
	r = middlewares.SetManuscriptID(r, newManuscriptID)
	app := middlewares.ApplicationFromRequest(r)

	ctx, err := app.SendCommand(r.Context(), commands.SubmitManuscript{
		Title:  dto.Title,
		Author: dto.Author,
		File:   dto.File,
		// FIXME: jamais utilisé ??
		FileName: dto.FileName,
	})
	if err != nil {
		slog.Error("manuscript creation request error", err)
		return html.RespondWithErrorFragment(w, r, "#new-manuscript-form", ManuscriptCreationError{})
	}
	r = r.WithContext(ctx)
	slog.Info("manuscript created", "manuscript_id", newManuscriptID.String())

	return html.RespondWithTemplate(w, r, WriterManuscriptDto{
		Title: dto.Title,
	}, "manuscript-item", html.WithFiles("manuscript-item.gohtml"), html.WithFuncs{
		Funcs: transateStatusFuncMap,
	})
}

type ManuscriptDto struct {
	Status string  `json:"status"`
	Url    url.URL `json:"url"`
}

func handleCancelManuscriptSubmission(w http.ResponseWriter, r *http.Request) *http.Request {
	manuscriptID := r.Context().Value(contexts.ManuscriptIDContextKey{}).(application.ManuscriptID)
	slog.Info("manuscript submission cancelling request", "manuscript_id", manuscriptID.String())
	app := middlewares.ApplicationFromRequest(r)
	ctx, err := app.SendCommand(r.Context(), commands.CancelManuscriptSubmission{})
	if err != nil {
		slog.Error("manuscript submission cancelling request error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageErrorAsJson(w, err)
		return r
	}
	r = r.WithContext(ctx)
	slog.Info("manuscript submission cancelled", "manuscript_id", manuscriptID)
	helpers.WriteJson(w, "")
	return r
}

func handleManuscriptState(w http.ResponseWriter, r *http.Request) *http.Request {
	manuscriptID := middlewares.GetManuscriptID(r)
	slog.Info("manuscript status request", "manuscript_id", manuscriptID.String())
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.Query(r.Context(), queries.ManuscriptState{})
	if err != nil {
		slog.Error("manuscript status query error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageErrorAsJson(w, err)
		return r
	}
	manuscript, castedSuccessfuly := queryResult.(domain.Manuscript)
	if !castedSuccessfuly {
		slog.Error("manuscript status query result casting error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageErrorAsJson(w, err)
		return r
	}

	helpers.WriteJson(w, ManuscriptDto{
		Status: string(manuscript.Status),
		Url:    manuscript.FileURL,
	})
	return r
}

func handleManuscriptsFuncs(
	serveMux *http.ServeMux,
	app application.Application,
	usersHistory application.UsersHistory,
	publicationsHistory application.PublicationsHistory,
	manuscriptsHistory application.ManuscriptsHistory,
	filesSaver ports.FilesSaver,
	jwtMiddleware middlewares.Middleware,
) {
	routes := []router.Route{
		{
			Path:   "/api/manuscripts",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.PersistNewEvents,
				middlewares.InjectContextualizedManuscriptsHistory,
				middlewares.EnsureUserIsAuthenticatedAndExtractUserID,
				middlewares.InjectFilesSaver(filesSaver),
				middlewares.InjectManuscriptsHistory(manuscriptsHistory),
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.InjectApplication(app),
			},
			Handler: handleManuscriptCreation,
		},
		{
			Path:   "/manuscripts",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.InjectContextualizedManuscriptsHistory,
				middlewares.EnsureUserIsAuthenticatedAndExtractUserID,
				middlewares.InjectManuscriptsHistory(manuscriptsHistory),
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.InjectApplication(app),
			},
			Handler: handleGetManuscripts,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.InjectContextualizedManuscriptsHistory,
				middlewares.InjectContextualizedUserHistory,
				middlewares.EnsureUserIsAuthenticatedAndExtractUserID,
				middlewares.ExtractManuscriptID,
				middlewares.InjectManuscriptsHistory(manuscriptsHistory),
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.InjectApplication(app),
				jwtMiddleware,
			},
			Handler: handleManuscriptState,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID/cancel-submission",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.PersistNewEvents,
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.InjectContextualizedManuscriptsHistory,
				middlewares.InjectContextualizedUserHistory,
				middlewares.EnsureUserIsAuthenticatedAndExtractUserID,
				middlewares.ExtractManuscriptID,
				middlewares.InjectManuscriptsHistory(manuscriptsHistory),
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.InjectApplication(app),
				jwtMiddleware,
			},
			Handler: handleCancelManuscriptSubmission,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID/review",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.PersistNewEvents,
				middlewares.OnlyAvailableForEditor,
				middlewares.InjectContextualizedManuscriptsHistory,
				middlewares.InjectContextualizedUserHistory,
				middlewares.EnsureUserIsAuthenticatedAndExtractUserID,
				middlewares.ExtractManuscriptID,
				middlewares.InjectPublicationsHistory(publicationsHistory),
				middlewares.InjectManuscriptsHistory(manuscriptsHistory),
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.InjectApplication(app),
				jwtMiddleware,
			},
			Handler: handleManuscriptReviewSubmission,
		},
		{
			Path:   "/manuscripts/to-review",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.OnlyAvailableForEditor,
				middlewares.InjectContextualizedManuscriptsHistory,
				middlewares.InjectContextualizedUserHistory,
				middlewares.EnsureUserIsAuthenticatedAndExtractUserID,
				middlewares.InjectManuscriptsHistory(manuscriptsHistory),
				middlewares.InjectUsersHistory(usersHistory),
				middlewares.InjectApplication(app),
			},
			Handler: handleGetManuscriptsToReview,
		},
	}
	router.HandleRoutes(serveMux, routes)
}
