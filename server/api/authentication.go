package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/helpers/auth0"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/api/router"
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

type Authenticator struct {
	*oicd.Provider
	oauth2.Config
}

func (a Authenticator) AuthCodeURL(state string) string {
	return "TODO"
}

func NewAuthenticator() (Authenticator, error) {
	domainUrl, err := url.Parse("https://" + auth0.Auth0Domain + "/")
	if err != nil {
		return Authenticator{}, fmt.Errorf("Failed to parse the issuer url: %v", err)
	}

	provider, err := oidc.NewProvider(
		context.Background(),
		domainUrl,
	)
	if err != nil {
		return Authenticator{}, err
	}

	conf := oauth2.Config{
		ClientID:     auth0.Auth0ClientId,
		ClientSecret: auth0.Auth0ClientSecret,
		RedirectURL:  auth0.Auth0CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

func handleLogin(authenticator Authenticator) func(http.ResponseWriter, *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		state, err := generateRandomState()
		if err != nil {
			slog.Error("login state generation has failed", err)
			helpers.ManageError(w, err)
			return r
		}

		stateCookie := http.Cookie{
			Name:     "state_cookie",
			Path:     "/",
			Value:    state,
			MaxAge:   int(time.Hour.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &stateCookie)
		slog.Info("redirecting to login page")
		http.Redirect(w, r, authenticator.AuthCodeURL(state), http.StatusFound)
		return r
	}
}

func handleAuthenticationFuncs(
	serveMux *http.ServeMux,
	app application.Application,
	authenticator Authenticator,
	jwtMiddleware middlewares.Middleware) {
	routes := []router.Route{
		{
			Path:    "/login",
			Method:  "GET",
			Handler: handleLogin(authenticator),
		},
	}
	router.HandleRoutes(serveMux, routes)
}
