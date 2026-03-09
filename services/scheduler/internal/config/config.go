package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretPort string
	VaultPort  string
	AuditPort  string
	DbUri      string
}

func Load() (*Config, error) {
	err := godotenv.Load(
		"../../.env",
		".env",
	)

	if err != nil {
		return &Config{}, err
	}

	return &Config{
		SecretPort: getEnv("SECRET_PORT", "50053"),
		VaultPort:  getEnv("VAULT_PORT", "50054"),
		AuditPort:  getEnv("AUDIT_PORT", "50055"),
		DbUri:      os.Getenv("DB_URI"),
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return ":" + val
	}

	return fallback
}
