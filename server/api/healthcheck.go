package api

import (
	"errors"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/persistency/mongodb"
	"golang.org/x/exp/slog"
)

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
			helpers.ManageErrorAsJson(w, errors.New("unhealthy database client"))
			return r
		}

		w.Write([]byte("ok"))
		return r
	}
}
