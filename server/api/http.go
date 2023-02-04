package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type SubmitManuscriptDto struct {
	ManuscriptName string `json:"manuscript_name"`
}

func handleManuscriptCreation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dto SubmitManuscriptDto
	err := decoder.Decode(&dto)
	slog.Info("receiving manuscript creation request", "body", dto)
	if err != nil {
		slog.Error("manuscript creation request dto decoding error", err)
		manageError(&w, err)
		return
	}

	returnedEvents, err := app.Send(commands.SubmitManuscript{
		ManuscriptName: dto.ManuscriptName,
	})
	if err != nil {
		slog.Error("manuscript creation request error", err)
		manageError(&w, err)
		return
	}
	// TODO: Générer l'id dans la partie infra et l'envoyer dans app.Send ?
	var manuscriptID events.ManuscriptID
	for _, nextEvent := range returnedEvents {
		manuscriptSubmittedEvent, foundExpectedEvent := nextEvent.(events.ManuscriptSubmitted)
		if foundExpectedEvent {
			manuscriptID = manuscriptSubmittedEvent.CreatedManuscriptID
		}
	}
	slog.Info("manuscript created", "manuscript_id", manuscriptID)
	w.Header().Add("Content-Type", "application/json")
	w.Write(
		[]byte(
			fmt.Sprintf("{\"id\": \"%v\"}", manuscriptID),
		),
	)
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

func handleGetManuscriptState(manuscriptID events.ManuscriptID, w http.ResponseWriter, r *http.Request) {
	slog.Info("manuscript status request", "manuscript_id", manuscriptID)
	queryResult, err := app.Query(queries.ManuscriptStatus{
		ManuscriptID: manuscriptID,
	})
	if err != nil {
		slog.Error("manuscript status query error", err)
		manageError(&w, err)
		return
	}
	status, castedSuccessfuly := queryResult.(domain.Status)
	if !castedSuccessfuly {
		slog.Error("manuscript status query result casting error", err)
		manageError(&w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	manuscriptJSON, err := json.Marshal(
		ManuscriptDto{
			Status: string(status),
		},
	)
	if err != nil {
		manageError(&w, err)
		return
	}
	w.Write(manuscriptJSON)
}

func handleCancelManuscriptSubmission(manuscriptID events.ManuscriptID, w http.ResponseWriter, r *http.Request) {
	slog.Info("manuscript submission cancelling request", "manuscript_id", manuscriptID)
	_, err := app.Send(commands.CancelManuscriptSubmission{
		ManuscriptID: manuscriptID,
	})
	if err != nil {
		slog.Error("manuscript submission cancelling request error", err)
		manageError(&w, err)
		return
	}
	slog.Info("manuscript submission cancelled", "manuscript_id", manuscriptID)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(""))
}

func handleManuscript(w http.ResponseWriter, r *http.Request) {
	// TODO: Dev un routing plus user friendly ou en utiliser un déjà dispo
	urlParts := strings.Split(r.URL.String(), "/")
	manuscriptID := urlParts[3]

	if r.Method == "GET" && len(urlParts) == 4 {
		handleGetManuscriptState(events.ManuscriptID(manuscriptID), w, r)
		return
	}

	if r.Method == "POST" && len(urlParts) == 5 && urlParts[4] == "cancel-submission" {
		handleCancelManuscriptSubmission(events.ManuscriptID(manuscriptID), w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "invalid_http_method")
}

func manageError(w *http.ResponseWriter, err error) {
	(*w).WriteHeader(http.StatusInternalServerError)
	(*w).Write([]byte(err.Error()))
}

var app application.Application

func Start() {
	slog.Info("start new application")
	app = application.NewApplication()
	slog.Info("setup HTTP API")

	http.HandleFunc("/api/manuscripts", handleManuscripts)
	http.HandleFunc("/api/manuscripts/", handleManuscript)

	// TODO: Variabiliser le port
	slog.Info("HTTP API start listening")
	http.ListenAndServe(":8080", nil)
}
