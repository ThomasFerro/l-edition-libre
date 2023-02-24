package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"

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

func authentifyAs(ctx context.Context, displayedName string) (context.Context, error) {
	var newUser api.CreateAccountResponseDto
	ctx, err := helpers.Call(ctx, helpers.HttpRequest{
		Url:    "http://localhost:8080/api/users",
		Method: http.MethodPost,
		Body: api.CreateAccountRequestDto{
			DisplayedName: displayedName,
		},
		ResponseDto: &newUser,
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
