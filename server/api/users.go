package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"golang.org/x/exp/slog"
)

type CreateAccountRequestDto struct {
	DisplayedName string `json:"displayed_name"`
}

type CreateAccountResponseDto struct {
	Id string `json:"id"`
}

func handleAccountCreation(w http.ResponseWriter, r *http.Request) *http.Request {
	decoder := json.NewDecoder(r.Body)
	var dto CreateAccountRequestDto
	err := decoder.Decode(&dto)
	slog.Info("receiving account creation request", "body", dto)
	if err != nil {
		slog.Error("account creation request dto decoding error", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}

	app := middlewares.ApplicationFromRequest(r)
	ctx, err := app.SendCommand(r.Context(), commands.CreateAccount{
		DisplayedName: dto.DisplayedName,
	})
	if err != nil {
		slog.Error("account creation request error", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}
	r = r.WithContext(ctx)
	newUserID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
	slog.Info("acount created", "user_id", newUserID)
	helpers.WriteJson(w, CreateAccountResponseDto{
		Id: string(newUserID),
	})
	return r
}

func handlePromoteToEditor(w http.ResponseWriter, r *http.Request) *http.Request {
	slog.Info("receiving promotion to editor request")

	// TODO: Plus proprement, via un middleware qui recup le param de l'url
	splitUrl := strings.Split(r.URL.Path, "/")
	userId := splitUrl[len(splitUrl)-1]
	r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey{}, application.UserID(userId)))

	app := middlewares.ApplicationFromRequest(r)
	ctx, err := app.SendCommand(r.Context(), commands.PromoteUserToEditor{})
	if err != nil {
		slog.Error("promotion to editor request error", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}
	r = r.WithContext(ctx)
	slog.Info("user promoted to editor")
	helpers.WriteJson(w, "")
	return r
}

func handleUsersFuncs(
	serveMux *http.ServeMux,
	app application.Application,
	userHistory application.UsersHistory,
	jwtMiddleware middlewares.Middleware) {
	routes := []router.Route{
		{
			Path:   "/api/users",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.PersistNewEvents,
				middlewares.InjectContextualizedUserHistory,
				middlewares.InjectUsersHistory(userHistory),
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
				jwtMiddleware,
			},
			Handler: handleAccountCreation,
		},
		{
			Path:   "/api/users/:userID",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.PersistNewEvents,
				middlewares.RequiresAdminApiKey,
				middlewares.InjectContextualizedUserHistory,
				middlewares.InjectUsersHistory(userHistory),
				middlewares.InjectApplication(app),
			},
			Handler: handlePromoteToEditor,
		},
	}
	router.HandleRoutes(serveMux, routes)
}
