package config

import "os"

// getenv returns environment variable by name or default value
func getenv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

var (
	Environment = struct {
		ServerPort        string
		Environment       string
		ServerReadTimeout string
		DatabaseServerURL string
	}{
		ServerPort:        getenv("CLOUD_API_SERVER_PORT", "8082"),
		Environment:       getenv("ENVIRONMENT", "development"),
		ServerReadTimeout: getenv("SERVER_READ_TIMEOUT", "60"),
		DatabaseServerURL: getenv("DB_SERVER_URL", "postgres://postgres:secret@localhost:5432/packsizer"),
	}
)
