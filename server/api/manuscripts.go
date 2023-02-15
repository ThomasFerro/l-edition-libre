package api

import (
	"encoding/json"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type WriterManuscriptDto struct {
	Title string `json:"title"`
}
type WriterManuscriptsDto struct {
	Manuscripts []WriterManuscriptDto `json:"manuscripts"`
}

func handleGetManuscripts(w http.ResponseWriter, r *http.Request) {
	slog.Info("writer manuscripts request")
	app := middlewares.ApplicationFromRequest(r)
	userID := middlewares.UserIdFromRequest(r)
	queryResult, err := app.ManuscriptsQuery(userID, queries.WriterManuscripts{})
	if err != nil {
		slog.Error("writer manuscripts query error", err)
		helpers.ManageError(w, err)
		return
	}
	manuscripts, castedSuccessfuly := queryResult.([]domain.Manuscript)
	if !castedSuccessfuly {
		slog.Error("writer manuscripts query result casting error", err)
		helpers.ManageError(w, err)
		return
	}

	dto := WriterManuscriptsDto{
		Manuscripts: []WriterManuscriptDto{},
	}
	for _, manuscript := range manuscripts {
		dto.Manuscripts = append(dto.Manuscripts, WriterManuscriptDto{
			Title: manuscript.Title,
		})
	}
	helpers.WriteJson(w, dto)
}

type SubmitManuscriptRequestDto struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type SubmitManuscriptResponseDto struct {
	Id string `json:"id"`
}

func handleManuscriptCreation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dto SubmitManuscriptRequestDto
	err := decoder.Decode(&dto)
	slog.Info("receiving manuscript creation request", "body", dto)
	if err != nil {
		slog.Error("manuscript creation request dto decoding error", err)
		helpers.ManageError(w, err)
		return
	}

	newManuscriptID := application.NewManuscriptID()
	r = middlewares.SetManuscriptID(r, newManuscriptID)
	app := middlewares.ApplicationFromRequest(r)
	_, err = app.SendCommand(r.Context(), commands.SubmitManuscript{
		Title:  dto.Title,
		Author: dto.Author,
	})
	if err != nil {
		slog.Error("manuscript creation request error", err)
		helpers.ManageError(w, err)
		return
	}
	slog.Info("manuscript created", "manuscript_id", newManuscriptID.String())
	helpers.WriteJson(w, SubmitManuscriptResponseDto{
		Id: newManuscriptID.String(),
	})
}

type ManuscriptDto struct {
	Status string `json:"status"`
}

func handleCancelManuscriptSubmission(w http.ResponseWriter, r *http.Request) {
	manuscriptID := r.Context().Value(contexts.ManuscriptIDContextKey).(application.ManuscriptID)
	// TODO: slog avec le context pour ne pas avoir à remettre les params chaque fois ?
	slog.Info("manuscript submission cancelling request", "manuscript_id", manuscriptID.String())
	app := middlewares.ApplicationFromRequest(r)
	_, err := app.SendCommand(r.Context(), commands.CancelManuscriptSubmission{})
	if err != nil {
		slog.Error("manuscript submission cancelling request error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageError(w, err)
		return
	}
	slog.Info("manuscript submission cancelled", "manuscript_id", manuscriptID)
	helpers.WriteJson(w, "")
}

func handleManuscriptState(w http.ResponseWriter, r *http.Request) {
	manuscriptID := middlewares.GetManuscriptID(r)
	slog.Info("manuscript status request", "manuscript_id", manuscriptID.String())
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.ManuscriptQuery(manuscriptID, queries.ManuscriptStatus{})
	if err != nil {
		slog.Error("manuscript status query error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageError(w, err)
		return
	}
	status, castedSuccessfuly := queryResult.(domain.Status)
	if !castedSuccessfuly {
		slog.Error("manuscript status query result casting error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageError(w, err)
		return
	}

	helpers.WriteJson(w, ManuscriptDto{
		Status: string(status),
	})
}

func handleManuscriptsFuncs(app application.Application) {
	routes := []router.Route{
		{
			Path:   "/api/manuscripts",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handleManuscriptCreation,
		},
		{
			Path:   "/api/manuscripts",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.InjectHistory(),
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handleGetManuscripts,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.InjectHistory(),
				middlewares.ExtractUserID,
				middlewares.ExtractManuscriptID,
				middlewares.InjectApplication(app),
			},
			Handler: handleManuscriptState,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID/cancel-submission",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.InjectHistory(),
				middlewares.ExtractUserID,
				middlewares.ExtractManuscriptID,
				middlewares.InjectApplication(app),
			},
			Handler: handleCancelManuscriptSubmission,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID/review",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.OnlyAvailableForEditor,
				middlewares.InjectHistory(),
				middlewares.ExtractUserID,
				middlewares.ExtractManuscriptID,
				middlewares.InjectApplication(app),
			},
			Handler: handleManuscriptReviewSubmission,
		},
		{
			Path:   "/api/manuscripts/to-review",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.OnlyAvailableForEditor,
				middlewares.InjectHistory(),
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handleGetManuscriptsToReview,
		},
	}
	router.HandleRoutes(routes)
}
