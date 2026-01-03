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

	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	return &Config{
		DBPort: dbPort,
		DBUser: dbUser,
		DBPass: dbPass,
		DBName: dbName,
	}
}
