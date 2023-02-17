package api

import (
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/configuration"
	"github.com/ThomasFerro/l-edition-libre/persistency"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

func Start() {
	slog.Info("start new application")
	managedCommands := application.ManagedCommands{
		"commands.CreateAccount":              commands.HandleCreateAccount,
		"commands.PromoteUserToEditor":        commands.HandlePromoteUserToEditor,
		"commands.SubmitManuscript":           commands.HandleSubmitManuscript,
		"commands.ReviewManuscript":           commands.HandleReviewManuscript,
		"commands.CancelManuscriptSubmission": commands.HandleCancelManuscriptSubmission,
	}
	managedQueries := application.ManagedQueries{
		"queries.ManuscriptStatus":    queries.HandleManuscriptStatus,
		"queries.WriterManuscripts":   queries.HandleWriterManuscripts,
		"queries.ManuscriptsToReview": queries.HandleManuscriptsToReview,
	}
	manuscriptsHistory := persistency.NewManuscriptsHistory()
	usersHistory := persistency.NewUsersHistory()
	app := application.NewApplication(managedCommands, managedQueries)
	slog.Info("setup HTTP API")

	handleManuscriptsFuncs(app, usersHistory, manuscriptsHistory)
	handleUsersFuncs(app, usersHistory)

	slog.Info("HTTP API start listening")
	port := configuration.GetConfiguration(configuration.PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
