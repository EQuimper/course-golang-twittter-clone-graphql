package config

import (
	"os"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type Config struct {
	Database database
}

func New() *Config {
	godotenv.Load()

	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
	}
}
