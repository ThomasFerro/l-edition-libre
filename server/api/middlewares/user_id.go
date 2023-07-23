package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"golang.org/x/exp/slog"
)

func TryGetUserIdFromRequest(r *http.Request) (application.UserID, bool) {
	value := r.Context().Value(contexts.UserIDContextKey{})
	if value == nil {
		return "", false
	}
	return value.(application.UserID), true
}

func UserIdFromRequest(r *http.Request) application.UserID {
	return r.Context().Value(contexts.UserIDContextKey{}).(application.UserID)
}

func ExtractUserID(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		userId, err := extractUserIDFromJwt(r)
		if err != nil {
			slog.Warn("Unable to extract user id: %v", err)
			http.Error(w, "Unable to extract user id", http.StatusBadRequest)
			return r
		}
		if userId == "" {
			return next(w, r)
		}

		r = r.WithContext(context.WithValue(r.Context(), contexts.UserIDContextKey{}, userId))

		return next(w, r)
	}
}

type UserInfoResponse struct {
	Sub string `json:"sub"`
}

func extractUserIDFromJwt(r *http.Request) (application.UserID, error) {
	token, err := extractBearerToken(r.Header)
	if err != nil {
		return "", err
	}

	// TODO: Extraire ? Récupérer de la logique du helper de tests ?
	userInfoUrl, err := url.Parse("https://" + auth0Domain + "/oauth/userinfo")
	if err != nil {
		return "", fmt.Errorf("Failed to parse the userinfo url: %v", err)
	}
	req, err := http.NewRequest("GET", userInfoUrl.String(), nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create the userinfo request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Userinfo request error: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("userinfo body read error: %v", err)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("wrong userinfo response code: %v %v", response.StatusCode, string(body))
	}
	var responseDto UserInfoResponse
	err = json.Unmarshal(body, &responseDto)
	if err != nil {
		return "", fmt.Errorf("userinfo body unmarshal error: %v (body: %v)", err, string(body))
	}
	return application.UserID(responseDto.Sub), nil
}

func extractBearerToken(header http.Header) (string, error) {
	authorizationHeader := header.Get("Authorization")
	split := strings.Split(authorizationHeader, "Bearer ")
	if len(split) != 2 {
		return "", errors.New("Unable to extract bearer")
	}
	return split[1], nil
}
