package api

import (
	"encoding/json"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

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
	app := middlewares.ApplicationFromRequest(r)
	_, err = app.SendManuscriptCommand(r.Context(), application.ManuscriptID(newManuscriptID), commands.SubmitManuscript{
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

func getManuscriptID(r *http.Request) application.ManuscriptID {
	return application.MustParseManuscriptID(helpers.FromUrlParams(r.Context(), ":manuscriptID"))
}

func handleCancelManuscriptSubmission(w http.ResponseWriter, r *http.Request) {
	manuscriptID := getManuscriptID(r)
	slog.Info("manuscript submission cancelling request", "manuscript_id", manuscriptID.String())
	app := middlewares.ApplicationFromRequest(r)
	_, err := app.SendManuscriptCommand(r.Context(), manuscriptID, commands.CancelManuscriptSubmission{})
	if err != nil {
		slog.Error("manuscript submission cancelling request error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageError(w, err)
		return
	}
	slog.Info("manuscript submission cancelled", "manuscript_id", manuscriptID)
	helpers.WriteJson(w, "")
}

func handleManuscriptState(w http.ResponseWriter, r *http.Request) {
	manuscriptID := getManuscriptID(r)
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
			Path:   "/api/manuscripts/:manuscriptID",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handleManuscriptState,
		},
		{
			Path:   "/api/manuscripts/:manuscriptID/cancel-submission",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handleCancelManuscriptSubmission,
		},
		// FIXME: Devrait être scoppé à l'éditeur
		{
			Path:   "/api/manuscripts/:manuscriptID/review",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.UserShouldHaveAccessToManuscript,
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handleManuscriptReviewSubmission,
		},
		{
			Path:   "/api/manuscripts/to-review",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.InjectApplication(app),
			},
			Handler: handleGetManuscriptsToReview,
		},
	}
	router.HandleRoutes(routes)
}
