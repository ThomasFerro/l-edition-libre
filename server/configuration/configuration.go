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
)

var ConfigurationsWithFallback = map[ConfigurationKey]string{
	PORT:                    "8080",
	ADMIN_API_KEY:           "ADMIN",
	MONGO_CONNECTION_STRING: "mongodb://127.0.0.1:27017/?directConnection=true",
	MONGO_DATABASE_NAME:     "l-edition-libre",
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
