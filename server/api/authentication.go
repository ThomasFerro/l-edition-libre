package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	*oidc.Provider
	oauth2.Config
}

func (a Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

func NewAuthenticator() (Authenticator, error) {
	domainUrl, err := url.Parse("https://" + auth0.Auth0Domain + "/")
	if err != nil {
		return Authenticator{}, fmt.Errorf("Failed to parse the issuer url: %v", err)
	}

	provider, err := oidc.NewProvider(
		context.Background(),
		domainUrl.String(),
	)
	if err != nil {
		return Authenticator{}, err
	}

	conf := oauth2.Config{
		ClientID:     auth0.Auth0ClientId,
		ClientSecret: auth0.Auth0ClientSecret,
		RedirectURL:  auth0.Auth0CallbackURL,
		Endpoint:     provider.Endpoint(),
		// TODO: Quels scopes ?
		Scopes: []string{oidc.ScopeOpenID, "profile"},
	}

	return Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

const STATE_COOKIE_NAME = "state_cookie"

func isAuthenticated(r *http.Request) (bool, error) {
	// TODO: Un appel pour vérifier si le token est ok ?
	stateCookie, err := r.Cookie(STATE_COOKIE_NAME)
	if errors.Is(err, http.ErrNoCookie) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stateCookie.Value != "", nil
}

func handleLogin(authenticator Authenticator) func(http.ResponseWriter, *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		state, err := generateRandomState()
		if err != nil {
			slog.Error("login state generation has failed", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}

		stateCookie := http.Cookie{
			Name:     STATE_COOKIE_NAME,
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
func handleLogout(w http.ResponseWriter, r *http.Request) *http.Request {
	logoutUrl, err := url.Parse("https://" + auth0.Auth0Domain + "/v2/logout")
	if err != nil {
		slog.Error("unable to parse logout url", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host + "/logout-callback")
	if err != nil {
		slog.Error("unable to parse callback url", err)
		helpers.ManageErrorAsJson(w, err)
		return r
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", auth0.Auth0ClientId)
	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
	return r
}

func handleLogoutCallback(w http.ResponseWriter, r *http.Request) *http.Request {
	stateTokenCookie := http.Cookie{
		Name:     "state_cookie",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &stateTokenCookie)
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &accessTokenCookie)
	profileCookie := http.Cookie{
		Name:     "profile",
		Path:     "/",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &profileCookie)
	http.Redirect(w, r, "/", http.StatusFound)
	return r
}

func handleLoginCallback(authenticator Authenticator) func(http.ResponseWriter, *http.Request) *http.Request {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		cookieState, err := r.Cookie("state_cookie")
		if err != nil {
			slog.Error("unable to get state cookie", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}
		if cookieState.Value != r.URL.Query().Get("state") {
			slog.Error("state mismatch")
			helpers.ManageErrorAsJson(w, errors.New("state mismatch"))
			return r
		}

		token, err := authenticator.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			slog.Error("unable to echange code for an auth token", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}

		idToken, err := authenticator.VerifyIDToken(r.Context(), token)
		if err != nil {
			slog.Error("unable to verify id token", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			slog.Error("unable to extract claims from id token", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}
		marshaledProfile, err := json.Marshal(profile)
		if err != nil {
			slog.Error("unable to marshal profile", err)
			helpers.ManageErrorAsJson(w, err)
			return r
		}

		accessTokenCookie := http.Cookie{
			Name:     "access_token",
			Path:     "/",
			Value:    token.AccessToken,
			MaxAge:   int(time.Hour.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &accessTokenCookie)
		profileCookie := http.Cookie{
			Name:     "profile",
			Path:     "/",
			Value:    string(marshaledProfile),
			MaxAge:   int(time.Hour.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &profileCookie)
		slog.Info("successfuly authenticated, redirecting to index")
		http.Redirect(w, r, "/", http.StatusFound)
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
		{
			Path:    "/callback",
			Method:  "GET",
			Handler: handleLoginCallback(authenticator),
		},
		{
			Path:    "/logout-callback",
			Method:  "GET",
			Handler: handleLogoutCallback,
		},
		{
			Path:    "/logout",
			Method:  "GET",
			Handler: handleLogout,
		},
	}
	router.HandleRoutes(serveMux, routes)
}
