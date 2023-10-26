package api

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
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
		authenticator := NewAuthenticator()

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

func handleIndexFuncs(serveMux *http.ServeMux) {
	routes := []router.Route{
		{
			Path:    "/",
			Method:  "GET",
			Handler: handleIndex(),
		},
	}
	router.HandleRoutes(serveMux, routes)
}

//go:embed html/index.go.html
var index string

type TemplateManuscriptDto struct {
	Name string
}
type IndexParameters struct {
	Manuscripts   []TemplateManuscriptDto
	Authenticated bool
}

func handleIndex() func(w http.ResponseWriter, r *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		fmt.Printf("cookies ? %v", r.Cookies())
		t, err := template.New("index").Parse(index)
		if err != nil {
			slog.Error("index template parsing error", err)
			helpers.ManageError(w, err)
			return r
		}
		err = t.Execute(w, IndexParameters{
			Manuscripts: []TemplateManuscriptDto{
				{
					Name: "Test",
				},
			},
		})
		if err != nil {
			slog.Error("index template execution error", err)
			helpers.ManageError(w, err)
			return r
		}
		return r
	}
}

func handleHealthCheckFuncs(serveMux *http.ServeMux, client *mongodb.DatabaseClient) {
	routes := []router.Route{
		{
			Path:    "/api/ready",
			Method:  "GET",
			Handler: handleApiIsReady(client),
		},
	}
	router.HandleRoutes(serveMux, routes)
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

func handleDatabaseFuncs(serveMux *http.ServeMux, client *mongodb.DatabaseClient) {
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
	router.HandleRoutes(serveMux, routes)
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
