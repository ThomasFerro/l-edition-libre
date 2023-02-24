package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func InjectFilesSaver(filesSaver commands.FilesSaver) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.FilesSaverContextKey{}, filesSaver))
			return next(w, r)
		}
	}
}
