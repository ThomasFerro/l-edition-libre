package auth0

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/configuration"
)

var Auth0ClientId = configuration.GetConfiguration(configuration.AUTH0_CLIENT_ID)
var Auth0ClientSecret = configuration.GetConfiguration(configuration.AUTH0_CLIENT_SECRET)
var Auth0CallbackURL = configuration.GetConfiguration(configuration.AUTH0_CALLBACK_URL)
var Auth0Domain = configuration.GetConfiguration(configuration.AUTH0_DOMAIN)
var Auth0Audience = configuration.GetConfiguration(configuration.AUTH0_AUDIENCE)

func ExtractUserID(token string) (application.UserID, error) {
	parts := strings.Split(token, ".")
	result, err := base64.RawURLEncoding.DecodeString(parts[1])

	if err != nil {
		return "", err
	}
	var payload Payload
	err = json.Unmarshal(result, &payload)
	if err != nil {
		return "", err
	}

	return application.UserID(payload.Sub), nil
}

type Payload struct {
	Sub string `json:"sub"`
}
