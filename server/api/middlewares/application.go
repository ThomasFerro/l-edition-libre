package middlewares

import (
	"context"
	"net/http"

	apiContext "github.com/ThomasFerro/l-edition-libre/api/context"
	"github.com/ThomasFerro/l-edition-libre/application"
)

func ApplicationFromRequest(r *http.Request) application.Application {
	return r.Context().Value(apiContext.ApplicationContextKey).(application.Application)
}

func InjectApplication(app application.Application, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), apiContext.ApplicationContextKey, app))
		next(w, r)
	}
}
