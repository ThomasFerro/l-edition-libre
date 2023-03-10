package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/configuration"
	"github.com/ThomasFerro/l-edition-libre/persistency/inmemory"
	"github.com/ThomasFerro/l-edition-libre/persistency/mongodb"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

func Start(databaseName string) {
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
	filesSaver := inmemory.NewFilesSaver()

	client, err := mongodb.GetClient(databaseName)
	if err != nil {
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

	handleHealthCheckFuncs(client)
	handleDatabaseFuncs(client)
	handleManuscriptsFuncs(app, usersHistory, publicationsHistory, manuscriptsHistory, filesSaver)
	handlePublicationsFuncs(app, publicationsHistory)
	handleUsersFuncs(app, usersHistory)

	slog.Info("HTTP API start listening")
	port := configuration.GetConfiguration(configuration.PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func handleHealthCheckFuncs(client *mongodb.DatabaseClient) {
	routes := []router.Route{
		{
			Path:    "/api/ready",
			Method:  "GET",
			Handler: handleApiIsReady(client),
		},
	}
	router.HandleRoutes(routes)
}

func handleApiIsReady(client *mongodb.DatabaseClient) func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		slog.Info("check api readiness")
		err := client.HealthCheck()
		if err != nil {
			slog.Error("unhealthy database client", err)
			helpers.ManageError(w, errors.New("unhealthy database client"))
			return r
		}

		w.Write([]byte("ok"))
		return r
	}
}

func handleDatabaseFuncs(client *mongodb.DatabaseClient) {
	routes := []router.Route{
		{
			Path:    "/api/init",
			Method:  "POST",
			Handler: handleInitDatabase(client),
			Middlewares: []middlewares.Middleware{
				middlewares.RequiresAdminApiKey,
			},
		},
	}
	router.HandleRoutes(routes)
}

func handleInitDatabase(client *mongodb.DatabaseClient) func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		slog.Info("initialize database")
		err := client.InitDatabase()
		if err != nil {
			slog.Error("unable to initialize database", err)
			helpers.ManageError(w, err)
			return r
		}

		w.Write([]byte("ok"))
		return r
	}
}
