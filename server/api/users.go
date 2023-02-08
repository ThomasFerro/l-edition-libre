package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
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
		manageError(w, err)
		return
	}

	newUserID := application.NewUserID()
	// TODO
	// _, err = app.Send(newUserID, commands.CreateAccount{
	// 	DisplayedName: dto.DisplayedName,
	// })
	if err != nil {
		slog.Error("account creation request error", err)
		manageError(w, err)
		return
	}
	slog.Info("acount created", "user_id", newUserID.String())
	writeJson(w, CreateAccountResponseDto{
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

func handleUsersFuncs() {
	http.HandleFunc("/api/users", handleUsers)
}

func extractUserId(r *http.Request) (application.UserID, error) {
	userId := r.Header.Get("X-User-Id")
	if userId == "" {
		return application.UserID{}, errors.New("user id not found")
	}
	return application.ParseUserID(userId)
}
