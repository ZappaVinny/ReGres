package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	appName              string
	corsAllowedOrigin    string
	httpsEnabled         bool
	sessionDurationHours int
	port                 string
)

func Load() {
	appName = getEnv("APP_NAME", "ReGres")
	corsAllowedOrigin = getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173")
	httpsEnabled = getEnvBool("HTTPS_ENABLED", false)
	sessionDurationHours = getEnvInt("SESSION_DURATION_HOURS", 168)
	port = getEnv("PORT", "8080")
}

func AppName() string {
	return appName
}

func SessionCookieName() string {
	return slugify(appName) + "_session"
}

func CorsAllowedOrigin() string {
	return corsAllowedOrigin
}

func HTTPSEnabled() bool {
	return httpsEnabled
}

func SessionDurationHours() int {
	return sessionDurationHours
}

func Port() string {
	return port
}

func slugify(value string) string {
	var b strings.Builder

	for _, r := range strings.ToLower(value) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
