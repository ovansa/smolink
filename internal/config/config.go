package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort      string
	PostgresDSN     string
	RedisAddr       string
	WebhookEndpoint string
	ClickThreshold  int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		ServerPort:      viper.GetString("SERVER_PORT"),
		PostgresDSN:     viper.GetString("POSTGRES_DSN"),
		RedisAddr:       viper.GetString("REDIS_ADDR"),
		WebhookEndpoint: viper.GetString("WEBHOOK_ENDPOINT"),
		ClickThreshold:  viper.GetInt("CLICK_THRESHOLD"),
	}

	return config, nil
}
