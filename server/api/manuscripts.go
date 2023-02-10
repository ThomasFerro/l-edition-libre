package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type SubmitManuscriptRequestDto struct {
	ManuscriptName string `json:"manuscript_name"`
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
		ManuscriptName: dto.ManuscriptName,
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

func handleManuscripts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleManuscriptCreation(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
	}
}

type ManuscriptDto struct {
	Status string `json:"status"`
}

func handleGetManuscriptState(manuscriptID application.ManuscriptID, w http.ResponseWriter, r *http.Request) {
	slog.Info("manuscript status request", "manuscript_id", manuscriptID.String())
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.Query(manuscriptID, queries.ManuscriptStatus{})
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

func handleCancelManuscriptSubmission(manuscriptID application.ManuscriptID, w http.ResponseWriter, r *http.Request) {
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

// type ReviewManuscriptRequestDto struct{}

func handleManuscriptReviewSubmission(manuscriptID application.ManuscriptID, w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// var dto ReviewManuscriptRequestDto
	// err := decoder.Decode(&dto)
	// slog.Info("manuscript submission review request", "user_id", userID.String(), "manuscript_id", manuscriptID.String(), "body", dto)
	// if err != nil {
	// 	slog.Error("manuscript creation request dto decoding error", err)
	// 	helpers.ManageError(w, err)
	// 	return
	// }
	app := middlewares.ApplicationFromRequest(r)
	_, err := app.SendManuscriptCommand(r.Context(), manuscriptID, commands.ReviewManuscript{})
	if err != nil {
		slog.Error("manuscript submission review request error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageError(w, err)
		return
	}
	slog.Info("manuscript submission reviewed", "manuscript_id", manuscriptID)
	helpers.WriteJson(w, "")
}

func handleManuscript(w http.ResponseWriter, r *http.Request) {
	// TODO: Dev un routing plus user friendly ou en utiliser un déjà dispo
	urlParts := strings.Split(r.URL.String(), "/")
	manuscriptID := r.Context().Value(contexts.ManuscriptIDContextKey).(application.ManuscriptID)

	if r.Method == "GET" && len(urlParts) == 4 {
		handleGetManuscriptState(application.ManuscriptID(manuscriptID), w, r)
		return
	}

	if r.Method == "POST" && len(urlParts) == 5 && urlParts[4] == "cancel-submission" {
		handleCancelManuscriptSubmission(application.ManuscriptID(manuscriptID), w, r)
		return
	}

	if r.Method == "POST" && len(urlParts) == 5 && urlParts[4] == "review" {
		handleManuscriptReviewSubmission(application.ManuscriptID(manuscriptID), w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "invalid_http_method")
}

func handleManuscriptsFuncs(app application.Application) {
	http.HandleFunc("/api/manuscripts", middlewares.InjectApplication(app, middlewares.ExtractUserID(handleManuscripts)))
	http.HandleFunc("/api/manuscripts/",
		middlewares.InjectApplication(app, middlewares.ExtractUserID(
			middlewares.ExtractManuscriptId(
				middlewares.UserShouldHaveAccessToManuscript(
					handleManuscript)))))
}
