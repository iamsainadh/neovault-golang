package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort string
	DB_URI   string
}

func LoadConfig() *Config {
	err := godotenv.Load(
		"../../.env",
		".env",
	)

	if err != nil {
		return &Config{}
	}
	return &Config{
		GRPCPort: getEnv("GRPCPort", "50054"),
		DB_URI:   mustGetEnv("DB_URI"),
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
