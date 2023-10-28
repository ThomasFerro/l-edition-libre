package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/configuration"
	"github.com/ThomasFerro/l-edition-libre/persistency/inmemory"
	"github.com/ThomasFerro/l-edition-libre/persistency/mongodb"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

func Start(databaseName string) *http.Server {
	slog.Info("start new application")
	port := configuration.GetConfiguration(configuration.PORT)
	serveMux := http.NewServeMux()
	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	go func() {
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
		filesSaver := inmemory.NewFilesSaver()

		client, err := mongodb.GetClient(databaseName)
		if err != nil {
			slog.Error("unable to get mongo client", err)
			return
		}
		defer func() {
			client.Close()
		}()
		usersHistory := mongodb.NewUsersHistory(client)
		publicationsHistory := mongodb.NewPublicationsHistory(client)
		manuscriptsHistory := mongodb.NewManuscriptsHistory(client)

		app := application.NewApplication(managedCommands, managedEvents, managedQueries)
		slog.Info("setup HTTP API")
		jwtMiddleware, err := middlewares.EnsureTokenIsValid()
		if err != nil {
			slog.Error("unable to get jwt middleware", err)
			return
		}

		authenticator, err := NewAuthenticator()
		if err != nil {
			slog.Error("unable to initiate authenticator", err)
			return
		}

		handleHealthCheckFuncs(serveMux, client)
		handleDatabaseFuncs(serveMux, client)
		handleIndexFuncs(serveMux)
		handleManuscriptsFuncs(serveMux, app, usersHistory, publicationsHistory, manuscriptsHistory, filesSaver, jwtMiddleware)
		handlePublicationsFuncs(serveMux, app, publicationsHistory, jwtMiddleware)
		handleUsersFuncs(serveMux, app, usersHistory, jwtMiddleware)
		handleAuthenticationFuncs(serveMux, app, authenticator, jwtMiddleware)

		slog.Info("HTTP API start listening")
		err = server.ListenAndServe()
		if err != nil {
			slog.Error("server listening error", err)
		}
	}()
	return server
}
