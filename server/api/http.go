package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/persistency"
	"golang.org/x/exp/slog"
)

func Start() {
	slog.Info("start new application")
	app := application.NewApplication(persistency.NewManuscriptsHistory(), persistency.NewUsersHistory())
	slog.Info("setup HTTP API")

	handleManuscriptsFuncs(app)
	handleUsersFuncs(app)

	// TODO: Variabiliser le port
	slog.Info("HTTP API start listening")
	http.ListenAndServe(":8080", nil)
}
