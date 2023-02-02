package testContext

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/go-bdd/gobdd/context"
)

type AppKey struct{}

func GetApp(ctx context.Context) *application.Application {
	app, err := ctx.Get(AppKey{})
	if err != nil {
		panic(err)
	}

	return app.(*application.Application)
}
