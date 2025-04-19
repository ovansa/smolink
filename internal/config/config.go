package config

import (
	"errors"

	"github.com/spf13/viper"
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
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("SERVER_PORT", ":8080")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{
		Environment:     viper.GetString("ENVIRONMENT"),
		ServerPort:      viper.GetString("SERVER_PORT"),
		PostgresDSN:     viper.GetString("POSTGRES_DSN"),
		RedisAddr:       viper.GetString("REDIS_ADDR"),
		RedisPassword:   viper.GetString("REDIS_PASSWORD"),
		RedisDB:         viper.GetInt("REDIS_DB"),
		WebhookEndpoint: viper.GetString("WEBHOOK_ENDPOINT"),
		ClickThreshold:  viper.GetInt("CLICK_THRESHOLD"),
	}

	if config.PostgresDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required")
	}

	return config, nil
}
