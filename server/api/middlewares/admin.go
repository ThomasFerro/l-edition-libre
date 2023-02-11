package middlewares

import (
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/configuration"
	"golang.org/x/exp/slog"
)

const ApiKeyHeader = "X-API-Key"

func RequiresAdminApiKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(ApiKeyHeader)
		configuredApiKey := configuration.GetConfiguration(configuration.ADMIN_API_KEY)
		if configuredApiKey != apiKey {
			slog.Warn("unauthorized api key")
			helpers.ManageError(w, fmt.Errorf("%v", http.StatusUnauthorized))
			return
		}
		next(w, r)
	}
}
