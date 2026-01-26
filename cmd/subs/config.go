package main

import (
	"log"

	"github.com/spf13/viper"
)

// Config хранит конфигурацию приложения
type Config struct {
	Env         string
	ServiceName string
	Port        string
}

// LoadConfig загружает конфигурацию
func LoadConfig() *Config {
	v := viper.New()
	v.SetConfigFile(".env") // можно использовать .env или config.yaml
	v.AutomaticEnv()        // fallback на env vars

	if err := v.ReadInConfig(); err != nil {
		log.Printf("onfig file not found, using env vars only: %v", err)
	}

	cfg := &Config{
		Env:         v.GetString("ENV"),
		ServiceName: v.GetString("SERVICE_NAME"),
		Port:        v.GetString("PORT"),
	}

	// базовая валидация
	if cfg.Env == "" || cfg.ServiceName == "" || cfg.Port == "" {
		log.Fatalf("ENV, SERVICE_NAME and PORT must be set")
	}

	return cfg
}
