package configuration

import (
	"os"

	"golang.org/x/exp/slog"
)

type ConfigurationKey string

const (
	PORT                    ConfigurationKey = "PORT"
	ADMIN_API_KEY           ConfigurationKey = "ADMIN_API_KEY"
	MONGO_CONNECTION_STRING ConfigurationKey = "MONGO_CONNECTION_STRING"
	MONGO_DATABASE_NAME     ConfigurationKey = "MONGO_DATABASE_NAME"
	AUTH0_DOMAIN            ConfigurationKey = "AUTH0_DOMAIN"
	AUTH0_AUDIENCE          ConfigurationKey = "AUTH0_AUDIENCE"
	AUTH0_CLIENT_ID         ConfigurationKey = "AUTH0_CLIENT_ID"
	AUTH0_CLIENT_SECRET     ConfigurationKey = "AUTH0_CLIENT_SECRET"
	AUTH0_CALLBACK_URL      ConfigurationKey = "AUTH0_CALLBACK_URL"
	AUTH0_USERINFO_URL      ConfigurationKey = "AUTH0_USERINFO_URL"
)

var ConfigurationsWithFallback = map[ConfigurationKey]string{
	PORT:                    "8080",
	ADMIN_API_KEY:           "ADMIN",
	MONGO_CONNECTION_STRING: "mongodb://127.0.0.1:27017/?directConnection=true",
	MONGO_DATABASE_NAME:     "l-edition-libre",
	AUTH0_DOMAIN:            "l-edition-libre.eu.auth0.com",
	AUTH0_AUDIENCE:          "https://leditionlibre/api",
	AUTH0_CLIENT_ID:         "xaUrOrVoyvMzgrEVPoNVO17eN4tCBOPR",
	AUTH0_CALLBACK_URL:      "http://localhost:8080/callback",
	AUTH0_USERINFO_URL:      "https://l-edition-libre.eu.auth0.com/userinfo",
}

func GetConfiguration(configurationKey ConfigurationKey) string {
	configurationFromEnv := os.Getenv(string(configurationKey))
	if configurationFromEnv != "" {
		return configurationFromEnv
	}
	fallback := ConfigurationsWithFallback[configurationKey]
	slog.Info(
		"No configuration found in env variables, falling back",
		"configuration_key", configurationKey,
		"fallback", fallback,
	)
	return fallback
}
