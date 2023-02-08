package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/persistency"
	"golang.org/x/exp/slog"
)

var app application.Application

func Start() {
	slog.Info("start new application")
	app = application.NewApplication(persistency.NewManuscriptsHistory())
	slog.Info("setup HTTP API")

	handleManuscriptsFuncs()
	handleUsersFuncs()

	// TODO: Variabiliser le port
	slog.Info("HTTP API start listening")
	http.ListenAndServe(":8080", nil)
}
