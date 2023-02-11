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

func handleAccountCreation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dto CreateAccountRequestDto
	err := decoder.Decode(&dto)
	slog.Info("receiving account creation request", "body", dto)
	if err != nil {
		slog.Error("account creation request dto decoding error", err)
		helpers.ManageError(w, err)
		return
	}

	newUserID := application.NewUserID()
	app := middlewares.ApplicationFromRequest(r)
	_, err = app.SendUserCommand(newUserID, commands.CreateAccount{
		DisplayedName: dto.DisplayedName,
	})
	if err != nil {
		slog.Error("account creation request error", err)
		helpers.ManageError(w, err)
		return
	}
	slog.Info("acount created", "user_id", newUserID.String())
	helpers.WriteJson(w, CreateAccountResponseDto{
		Id: newUserID.String(),
	})
}

func handlePromoteToEditor(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.UserIdFromRequest(r)
	slog.Info("receiving promotion to editor request", "user_id", userID)

	app := middlewares.ApplicationFromRequest(r)
	_, err := app.SendUserCommand(userID, commands.PromoteUserToEditor{})
	if err != nil {
		slog.Error("promotion to editor request error", err)
		helpers.ManageError(w, err)
		return
	}
	slog.Info("user promoted to editor", "user_id", userID)
	helpers.WriteJson(w, "")
}

func handleUsersFuncs(app application.Application) {
	routes := []router.Route{
		{
			Path:   "/api/users",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.InjectApplication(app),
			},
			Handler: handleAccountCreation,
		},
		{
			Path:   "/api/users/:userID",
			Method: "POST",
			Middlewares: []middlewares.Middleware{
				middlewares.RequiresAdminApiKey,
				middlewares.ExtractUserID,
				middlewares.InjectApplication(app),
			},
			Handler: handlePromoteToEditor,
		},
	}
	router.HandleRoutes(routes)
}
