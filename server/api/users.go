package api

import (
	"encoding/json"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
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
		helpers.ManageError(w, err)
		return r
	}

	newUserID := application.NewUserID()
	r = middlewares.SetUserId(r, newUserID)
	app := middlewares.ApplicationFromRequest(r)
	ctx, err := app.SendCommand(r.Context(), commands.CreateAccount{
		DisplayedName: dto.DisplayedName,
	})
	if err != nil {
		slog.Error("account creation request error", err)
		helpers.ManageError(w, err)
		return r
	}
	r = r.WithContext(ctx)
	slog.Info("acount created", "user_id", newUserID.String())
	helpers.WriteJson(w, CreateAccountResponseDto{
		Id: newUserID.String(),
	})
	return r
}

func handlePromoteToEditor(w http.ResponseWriter, r *http.Request) *http.Request {
	slog.Info("receiving promotion to editor request")

	app := middlewares.ApplicationFromRequest(r)
	ctx, err := app.SendCommand(r.Context(), commands.PromoteUserToEditor{})
	if err != nil {
		slog.Error("promotion to editor request error", err)
		helpers.ManageError(w, err)
		return r
	}
	r = r.WithContext(ctx)
	slog.Info("user promoted to editor")
	helpers.WriteJson(w, "")
	return r
}

func handleUsersFuncs(
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
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
				jwtMiddleware,
			},
			Handler: handlePromoteToEditor,
		},
	}
	router.HandleRoutes(routes)
}
