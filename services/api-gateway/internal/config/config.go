package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT        string
	AUTH_PORT   string
	SECRET_PORT string
	VAULT_PORT  string
	AUDIT_PORT  string
	JwtSecret   string
	DB_URI      string
	KAFKA_URI   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(
		"../../.env",
		".env",
	)
	if err != nil {
		log.Println("No .env file found, using system environment variables")
		return &Config{}, err
	}
	return &Config{
		PORT:        ":" + os.Getenv("PORT"),
		AUTH_PORT:   ":" + os.Getenv("AUTH_PORT"),
		SECRET_PORT: ":" + os.Getenv("SECRET_PORT"),
		VAULT_PORT:  ":" + os.Getenv("VAULT_PORT"),
		AUDIT_PORT:  ":" + os.Getenv("AUDIT_PORT"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
		DB_URI:      os.Getenv("DB_URI"),
		KAFKA_URI:   os.Getenv("KAFKA_URI"),
	}, nil
}
