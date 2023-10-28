package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/persistency/mongodb"
	"golang.org/x/exp/slog"
)

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
