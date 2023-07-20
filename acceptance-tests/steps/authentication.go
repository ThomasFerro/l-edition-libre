package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/configuration"
	"github.com/cucumber/godog"
)

func authentifyAsWriter(ctx context.Context) (context.Context, error) {
	return authentifyAs(ctx, "Writer")
}

func authentifyAsWriterWithName(ctx context.Context, writer string) (context.Context, error) {
	return authentifyAs(ctx, writer)
}

func authentifyAsEditor(ctx context.Context) (context.Context, error) {
	ctx, err := authentifyAs(ctx, "Editor")
	if err != nil {
		return ctx, fmt.Errorf("unable to authentify: %v", err)
	}
	headers := make(map[string]string, 0)
	headers[middlewares.ApiKeyHeader] = configuration.GetConfiguration(configuration.ADMIN_API_KEY)
	return helpers.Call(
		ctx,
		helpers.HttpRequest{
			Url:     "http://localhost:8080/api/users/promote-to-editor",
			Method:  http.MethodPost,
			Headers: headers,
		})
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (t TokenResponse) String() string {
	return fmt.Sprintf("TokenResponse{ AccessToken: %v }\n", t.AccessToken)
}

func getTokenFor(ctx context.Context, displayedName string) (string, error) {
	// TODO: Switch entre les variables en fonction du display name
	path := configuration.GetConfiguration("AUTH0_PATH")
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("audience", configuration.GetConfiguration("AUTH0_AUDIENCE"))
	data.Set("username", configuration.GetConfiguration("AUTH0_WRITER_USERNAME"))
	data.Set("password", configuration.GetConfiguration("AUTH0_WRITER_PASSWORD"))
	data.Set("client_id", configuration.GetConfiguration("AUTH0_CLIENT_ID"))
	data.Set("client_secret", configuration.GetConfiguration("AUTH0_CLIENT_SECRET"))

	r, err := http.NewRequest(http.MethodPost, path, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var tokenResponse TokenResponse
	ctx, err = helpers.DoCall(ctx, r, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func authentifyAs(ctx context.Context, displayedName string) (context.Context, error) {
	token, err := getTokenFor(ctx, displayedName)
	if err != nil {
		return nil, err
	}

	var newUser api.CreateAccountResponseDto
	ctx, err = helpers.Call(ctx, helpers.HttpRequest{
		Url:    "http://localhost:8080/api/users",
		Method: http.MethodPost,
		Body: api.CreateAccountRequestDto{
			DisplayedName: displayedName,
		},
		ResponseDto: &newUser,
		Token:       token,
	})
	if err != nil {
		return ctx, fmt.Errorf("unable to create a new account: %v", err)
	}

	newUserID := application.MustParseUserID(newUser.Id)
	ctx = helpers.SetUserName(ctx, newUserID, displayedName)
	return helpers.SetAuthentifiedUserID(ctx, newUserID), nil
}

func AuthenticationSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I am an authentified editor$`, authentifyAsEditor)
	ctx.Step(`^I authentify as an editor$`, authentifyAsEditor)
	ctx.Step(`I am an authentified writer`, authentifyAsWriter)
	ctx.Step(`I am authentified as another writer`, authentifyAsWriter)
}
