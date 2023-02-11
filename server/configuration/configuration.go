package configuration

import (
	"os"

	"golang.org/x/exp/slog"
)

type ConfigurationKey string

const (
	PORT          ConfigurationKey = "PORT"
	ADMIN_API_KEY ConfigurationKey = "ADMIN_API_KEY"
)

var ConfigurationsWithFallback = map[ConfigurationKey]string{
	PORT:          "8080",
	ADMIN_API_KEY: "ADMIN",
}

func GetConfiguration(configurationKey ConfigurationKey) string {
	fallback := ConfigurationsWithFallback[configurationKey]
	configurationFromEnv := os.Getenv(string(configurationKey))
	if configurationFromEnv == "" {
		configurationFromEnv = fallback
		slog.Info(
			"No configuration found in env variables, falling back",
			"configuration_key", configurationKey,
			"fallback", fallback,
		)
	}
	return configurationFromEnv
}
