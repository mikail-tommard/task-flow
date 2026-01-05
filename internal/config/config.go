package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func New() *Config {
	_ = godotenv.Load()

	return &Config{
		DBPort: getEnv("DB_PORT", "5433"),
		DBUser: getEnv("DB_USER", "tasksflowuser"),
		DBPass: getEnv("DB_PASS", "tasksflowpass"),
		DBName: getEnv("DB_NAME", "taskflowdb"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return  value
	}
	return defaultValue
}