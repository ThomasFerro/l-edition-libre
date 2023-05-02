package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type PublicationDto struct {
	Status string `json:"status"`
}

func handlePublicationState(w http.ResponseWriter, r *http.Request) *http.Request {
	publicationID := middlewares.GetPublicationID(r)
	slog.Info("publication status request", "publication_id", publicationID.String())
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.Query(r.Context(), queries.PublicationStatus{})
	if err != nil {
		slog.Error("publication status query error", err, "publication_id", publicationID.String())
		helpers.ManageError(w, err)
		return r
	}
	status, castedSuccessfuly := queryResult.(domain.PublicationStatus)
	if !castedSuccessfuly {
		slog.Error("publication status query result casting error", err, "publication_id", publicationID.String())
		helpers.ManageError(w, err)
		return r
	}

	helpers.WriteJson(w, PublicationDto{
		Status: string(status),
	})
	return r
}

func handlePublicationsFuncs(
	app application.Application,
	publicationsHistory application.PublicationsHistory,
	jwtMiddleware middlewares.Middleware) {
	routes := []router.Route{
		{
			Path:   "/api/publications/:publicationID",
			Method: "GET",
			Middlewares: []middlewares.Middleware{
				middlewares.InjectContextualizedPublicationHistory,
				middlewares.ExtractPublicationID,
				middlewares.InjectPublicationsHistory(publicationsHistory),
				middlewares.InjectApplication(app),
				jwtMiddleware,
			},
			Handler: handlePublicationState,
		},
	}
	router.HandleRoutes(routes)
}
