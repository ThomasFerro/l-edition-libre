package auth0

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/configuration"
)

var Auth0ClientId = configuration.GetConfiguration(configuration.AUTH0_CLIENT_ID)
var Auth0ClientSecret = configuration.GetConfiguration(configuration.AUTH0_CLIENT_SECRET)
var Auth0CallbackURL = configuration.GetConfiguration(configuration.AUTH0_CALLBACK_URL)
var Auth0UserInfoURL = configuration.GetConfiguration(configuration.AUTH0_USERINFO_URL)
var Auth0Domain = configuration.GetConfiguration(configuration.AUTH0_DOMAIN)
var Auth0Audience = configuration.GetConfiguration(configuration.AUTH0_AUDIENCE)

func ExtractUserID(token string) (application.UserID, error) {
	request, err := http.NewRequest("GET", Auth0UserInfoURL, nil)
	if err != nil {
		return application.UserID(""), fmt.Errorf("unable to create user info request: %v", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return application.UserID(""), fmt.Errorf("unable to get user info: %v", err)
	}
	defer response.Body.Close()

	err = handleHttpError(response)
	if err != nil {
		return application.UserID(""), fmt.Errorf("user info error: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return application.UserID(""), fmt.Errorf("user info body reading error: %v (body: %v)", err, body)
	}
	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return application.UserID(""), fmt.Errorf("user info body unmarshal error: %v (body: %v)", err, string(body))
	}
	return application.UserID(payload.Sub), nil
}

func handleHttpError(response *http.Response) error {
	if response.StatusCode >= 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("body read error: %v", err)
		}
		return fmt.Errorf("wrong response code: %v (%v)", response.StatusCode, string(body))
	}
	return nil
}

type Payload struct {
	Sub string `json:"sub"`
}
