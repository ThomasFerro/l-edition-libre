package middlewares

import (
	"context"
	"net/http"
	"strings"

	apiContext "github.com/ThomasFerro/l-edition-libre/api/context"
	"github.com/ThomasFerro/l-edition-libre/application"
)

func ExtractManuscriptId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/manuscripts/") {
			next(w, r)
			return
		}
		urlParts := strings.Split(r.URL.String(), "/")
		manuscriptID := application.MustParseManuscriptID(urlParts[3])
		r = r.WithContext(context.WithValue(r.Context(), apiContext.ManuscriptIDContextKey, manuscriptID))
		next(w, r)
	}
}
