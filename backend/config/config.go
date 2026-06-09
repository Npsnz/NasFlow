package config

import (
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

var AppConfig Config

func LoadConfig() {
	AppConfig.Port = getEnv("PORT", "8080")
	AppConfig.JWTSecret = getEnv("JWT_SECRET", "taskflow-super-secret-key-1234567890abcdef")
	AppConfig.DBPath = getEnv("DB_PATH", "taskflow.db")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
