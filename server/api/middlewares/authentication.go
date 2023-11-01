package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ThomasFerro/l-edition-libre/api/helpers/auth0"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"golang.org/x/exp/slog"
)

type CustomClaims struct {
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func EnsureTokenIsValid() (Middleware, error) {
	issuerURL, err := url.Parse("https://" + auth0.Auth0Domain + "/")
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{auth0.Auth0Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to set up the jwt validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("Encountered error while validating JWT", err)

		// TODO: Remplacer par une déco + redirection sur l'index en indiquant qu'il faut s'authentifier à nouveau ?
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next HandlerFuncReturningRequest) HandlerFuncReturningRequest {
		return func(w http.ResponseWriter, r *http.Request) *http.Request {
			var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
				next(w, r)
			}

			middleware.CheckJWT(handler).ServeHTTP(w, r)

			return r
		}
	}, nil
}
