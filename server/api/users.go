package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
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

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleAccountCreation(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
	}
}

func handleUsersFuncs(app application.Application) {
	http.HandleFunc("/api/users", middlewares.InjectApplication(app, handleUsers))
}
