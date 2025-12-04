package configs

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Sentry   SentryConfig   `json:"sentry"`
}

type SentryConfig struct {
	DSN         string
	Environment string
	Release     string
}

type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Password string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	file, err := os.ReadFile("configs/config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Sentry.DSN = os.Getenv("SENTRY_DSN")
	config.Sentry.Environment = os.Getenv("SENTRY_ENVIRONMENT")
	if config.Sentry.Environment == "" {
		config.Sentry.Environment = "production"
	}
	config.Sentry.Release = os.Getenv("SENTRY_RELEASE")

	return &config, nil
}
