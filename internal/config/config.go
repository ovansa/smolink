package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment     string
	ServerPort      string
	PostgresDSN     string
	RedisAddr       string
	RedisDB         int
	RedisPassword   string
	WebhookEndpoint string
	ClickThreshold  int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	// Helper function to get env vars with fallback
	getEnv := func(key, fallback string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return fallback
	}

	// Helper to get integer env vars
	getEnvInt := func(key string, fallback int) int {
		if value, exists := os.LookupEnv(key); exists {
			if intValue, err := strconv.Atoi(value); err == nil {
				return intValue
			}
		}
		return fallback
	}

	port := getEnv("PORT", getEnv("SERVER_PORT", "8080"))
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	config := &Config{
		Environment:     getEnv("ENVIRONMENT", "development"),
		ServerPort:      port, // Render uses PORT
		PostgresDSN:     getEnv("POSTGRES_DSN", ""),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvInt("REDIS_DB", 0),
		WebhookEndpoint: getEnv("WEBHOOK_ENDPOINT", ""),
		ClickThreshold:  getEnvInt("CLICK_THRESHOLD", 10), // Default 10 clicks
	}

	// Validate required configuration
	if config.PostgresDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required")
	}

	if !strings.HasPrefix(config.ServerPort, ":") {
		return nil, fmt.Errorf("invalid SERVER_PORT: must begin with ':' (got %s)", config.ServerPort)
	}

	return config, nil
}
