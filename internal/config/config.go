package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
}

func Load() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		DatabaseURL: 	strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT is required")
	}

if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}


	return cfg, nil;

}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
    return fallback
	}

	return value
}