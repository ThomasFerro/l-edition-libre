package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/cucumber/godog"
)

func authentifyAsWriter(ctx context.Context) (context.Context, error) {
	return authentifyAs(ctx, "Writer")
}

func authentifyAsEditor(ctx context.Context) (context.Context, error) {
	// TODO: Ajouter le droit
	return authentifyAs(ctx, "Editor")
}

func authentifyAs(ctx context.Context, displayedName string) (context.Context, error) {
	var newUser api.CreateAccountResponseDto
	ctx, err := helpers.Call(ctx, "http://localhost:8080/api/users", http.MethodPost, api.CreateAccountRequestDto{
		DisplayedName: displayedName,
	}, &newUser)
	if err != nil {
		return ctx, fmt.Errorf("unable to create a new account: %v", err)
	}

	newUserID := application.MustParseUserID(newUser.Id)
	return helpers.SetAuthentifiedUserID(ctx, newUserID), nil
}

func AuthenticationSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I am an authentified editor$`, authentifyAsEditor)
	ctx.Step(`I am an authentified writer`, authentifyAsWriter)
	ctx.Step(`I am authentified as another writer`, authentifyAsWriter)
}
