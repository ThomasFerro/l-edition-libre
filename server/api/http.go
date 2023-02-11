package api

import (
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/configuration"
	"github.com/ThomasFerro/l-edition-libre/persistency"
	"golang.org/x/exp/slog"
)

func Start() {
	slog.Info("start new application")
	app := application.NewApplication(persistency.NewManuscriptsHistory(), persistency.NewUsersHistory())
	slog.Info("setup HTTP API")

	handleManuscriptsFuncs(app)
	handleUsersFuncs(app)

	slog.Info("HTTP API start listening")
	port := configuration.GetConfiguration(configuration.PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
