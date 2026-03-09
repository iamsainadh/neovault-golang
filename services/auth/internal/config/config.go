package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT       string
	DB_URI     string
	VAULT_PORT string
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
		PORT:       os.Getenv("PORT"),
		DB_URI:     os.Getenv("DB_URI"),
		VAULT_PORT: os.Getenv("VAULT_PORT"),
	}
}
