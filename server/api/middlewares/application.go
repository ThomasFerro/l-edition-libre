package middlewares

import (
	"context"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func ApplicationFromRequest(r *http.Request) application.Application {
	return r.Context().Value(contexts.ApplicationContextKey).(application.Application)
}

func InjectApplication(app application.Application) Middleware {
	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			r = r.WithContext(context.WithValue(r.Context(), contexts.ApplicationContextKey, app))
			return next(w, r)
		}
	}
}
