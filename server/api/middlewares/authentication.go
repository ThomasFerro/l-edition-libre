package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ThomasFerro/l-edition-libre/configuration"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"golang.org/x/exp/slog"
)

var auth0Domain = configuration.GetConfiguration(configuration.AUTH0_DOMAIN)
var auth0Audience = configuration.GetConfiguration(configuration.AUTH0_AUDIENCE)

type CustomClaims struct {
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func EnsureTokenIsValid() (Middleware, error) {
	issuerURL, err := url.Parse("https://" + auth0Domain + "/")
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{auth0Audience},
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

// func EnsureTokenIsValid(handler http.Handler) http.Handler {
// 	issuerURL, err := url.Parse("https://" + auth0Domain + "/")
// 	if err != nil {
// 		slog.Error("Failed to parse the issuer url", err)
// 		// TODO: Mieux qu'un panic ?
// 		panic("Failed to parse the issuer url")
// 	}

// 	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

// 	jwtValidator, err := validator.New(
// 		provider.KeyFunc,
// 		validator.RS256,
// 		issuerURL.String(),
// 		[]string{auth0Audience},
// 		validator.WithCustomClaims(
// 			func() validator.CustomClaims {
// 				return &CustomClaims{}
// 			},
// 		),
// 		validator.WithAllowedClockSkew(time.Minute),
// 	)
// 	if err != nil {
// 		slog.Error("Failed to set up the jwt validator")
// 		panic("Failed to set up the jwt validator")
// 	}

// 	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
// 		slog.Error("Encountered error while validating JWT", err)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
// 	}

// 	middleware := jwtmiddleware.New(
// 		jwtValidator.ValidateToken,
// 		jwtmiddleware.WithErrorHandler(errorHandler),
// 	)

// 	return middleware.CheckJWT(handler)
// }
