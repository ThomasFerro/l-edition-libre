package godogs_test

import (
	"acceptance-tests/helpers"
	"acceptance-tests/steps"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		databaseName := fmt.Sprintf("l-edition-libre-acceptance-%v", uuid.New().String())
		server := api.Start(databaseName)
		ctx = context.WithValue(ctx, "server", server)
		return context.WithValue(ctx, helpers.TagsKey{}, sc.Tags), nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		server := ctx.Value("server").(*http.Server)
		shutdownErr := server.Shutdown(ctx)
		if shutdownErr != nil {
			return nil, fmt.Errorf("unable to shutdown the server: %v", shutdownErr)
		}
		unhandledError, ok := ctx.Value(helpers.ErrorKey{}).(string)
		if ok {
			return ctx, fmt.Errorf("unhandled error in tests: %v", unhandledError)
		}
		return ctx, nil
	})

	steps.AuthenticationSteps(ctx)
	steps.ManuscriptSteps(ctx)
	steps.EditorSteps(ctx)
	steps.ErrorSteps(ctx)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
