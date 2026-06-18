package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisAddr   string
	ModelsDir   string
	GinMode     string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: buildDatabaseURL(),
		RedisAddr:   getEnv("REDIS_ADDR", "loaclhost:6379"),
		ModelsDir:   getEnv("MODELS_DIR", "./models"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}
}

func buildDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getEnv("DB_USER", "inference"),
		getEnv("DB_PASSWORD", "interdrence_secret"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "inference_db"),
	)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
