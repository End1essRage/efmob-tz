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
	PostgresDSN string // например "host=localhost user=postgres password=pass dbname=subs port=5432 sslmode=disable"
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
		PostgresDSN: v.GetString("POSTGRES_DSN"),
	}

	// базовая валидация
	if cfg.Env == "" || cfg.ServiceName == "" || cfg.Port == "" {
		log.Fatalf("ENV, SERVICE_NAME and PORT must be set")
	}

	return cfg
}
