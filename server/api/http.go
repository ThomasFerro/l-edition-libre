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
	managedEvents := application.ManagedEvents{
		"events.ManuscriptReviewed": application.HandleManuscriptReviewed,
	}
	managedQueries := application.ManagedQueries{
		"queries.ManuscriptState":     queries.HandleManuscriptState,
		"queries.WriterManuscripts":   queries.HandleWriterManuscripts,
		"queries.ManuscriptsToReview": queries.HandleManuscriptsToReview,
		"queries.PublicationStatus":   queries.HandlePublicationStatus,
	}
	manuscriptsHistory := persistency.NewManuscriptsHistory()
	publicationsHistory := persistency.NewPublicationsHistory()
	usersHistory := persistency.NewUsersHistory()
	filesSaver := persistency.NewFilesSaver()
	app := application.NewApplication(managedCommands, managedEvents, managedQueries)
	slog.Info("setup HTTP API")

	handleManuscriptsFuncs(app, usersHistory, publicationsHistory, manuscriptsHistory, filesSaver)
	handlePublicationsFuncs(app, publicationsHistory)
	handleUsersFuncs(app, usersHistory)

	slog.Info("HTTP API start listening")
	port := configuration.GetConfiguration(configuration.PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
