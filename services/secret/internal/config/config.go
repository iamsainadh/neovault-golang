package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort   string
	DB_URI     string
	VAULT_PORT string
}

func LoadConfig() *Config {
	godotenv.Load(
		"../../.env",
		".env",
	)

	return &Config{
		GRPCPort:   getEnv("GRPCPort", "50053"),
		DB_URI:     mustGetEnv("DB_URI"),
		VAULT_PORT: ":" + getEnv("VAULT_PORT", "50054"),
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
