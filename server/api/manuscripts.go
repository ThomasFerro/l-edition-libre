package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
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
	userId, err := extractUserId(r)
	if err != nil {
		manageError(w, err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var dto SubmitManuscriptRequestDto
	err = decoder.Decode(&dto)
	slog.Info("receiving manuscript creation request", "body", dto)
	if err != nil {
		slog.Error("manuscript creation request dto decoding error", err)
		manageError(w, err)
		return
	}

	newManuscriptID := application.NewManuscriptID()
	_, err = app.Send(userId, application.ManuscriptID(newManuscriptID), commands.SubmitManuscript{
		ManuscriptName: dto.ManuscriptName,
	})
	if err != nil {
		slog.Error("manuscript creation request error", err)
		manageError(w, err)
		return
	}
	slog.Info("manuscript created", "manuscript_id", newManuscriptID.String())
	writeJson(w, SubmitManuscriptResponseDto{
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

func handleGetManuscriptState(userID application.UserID, manuscriptID application.ManuscriptID, w http.ResponseWriter, r *http.Request) {
	slog.Info("manuscript status request", "user_id", userID.String(), "manuscript_id", manuscriptID.String())
	queryResult, err := app.Query(userID, manuscriptID, queries.ManuscriptStatus{})
	if err != nil {
		slog.Error("manuscript status query error", err)
		manageError(w, err)
		return
	}
	status, castedSuccessfuly := queryResult.(domain.Status)
	if !castedSuccessfuly {
		slog.Error("manuscript status query result casting error", err)
		manageError(w, err)
		return
	}

	writeJson(w, ManuscriptDto{
		Status: string(status),
	})
}

func handleCancelManuscriptSubmission(userID application.UserID, manuscriptID application.ManuscriptID, w http.ResponseWriter, r *http.Request) {
	slog.Info("manuscript submission cancelling request", "user_id", userID.String(), "manuscript_id", manuscriptID.String())
	_, err := app.Send(userID, manuscriptID, commands.CancelManuscriptSubmission{})
	if err != nil {
		slog.Error("manuscript submission cancelling request error", err)
		manageError(w, err)
		return
	}
	slog.Info("manuscript submission cancelled", "manuscript_id", manuscriptID)
	writeJson(w, "")
}

func handleManuscript(w http.ResponseWriter, r *http.Request) {
	userId, err := extractUserId(r)
	if err != nil {
		manageError(w, err)
		return
	}
	// TODO: Dev un routing plus user friendly ou en utiliser un déjà dispo
	urlParts := strings.Split(r.URL.String(), "/")
	manuscriptID := application.MustParseManuscriptID(urlParts[3])

	if r.Method == "GET" && len(urlParts) == 4 {
		handleGetManuscriptState(userId, application.ManuscriptID(manuscriptID), w, r)
		return
	}

	if r.Method == "POST" && len(urlParts) == 5 && urlParts[4] == "cancel-submission" {
		handleCancelManuscriptSubmission(userId, application.ManuscriptID(manuscriptID), w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "invalid_http_method")
}

func handleManuscriptsFuncs() {
	http.HandleFunc("/api/manuscripts", handleManuscripts)
	http.HandleFunc("/api/manuscripts/", handleManuscript)
}
