package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func ExtractManuscriptId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/manuscripts/") {
			next(w, r)
			return
		}
		urlParts := strings.Split(r.URL.String(), "/")
		manuscriptID := application.MustParseManuscriptID(urlParts[3])
		r = r.WithContext(context.WithValue(r.Context(), contexts.ManuscriptIDContextKey, manuscriptID))
		next(w, r)
	}
}
