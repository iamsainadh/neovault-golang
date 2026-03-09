package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	KafkaBrokerUri string
	DbUri          string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return &Config{
		Port:           ":" + os.Getenv("PORT"),
		KafkaBrokerUri: os.Getenv("KAFKA_URI"),
		DbUri:          os.Getenv("DB_uri"),
	}
}
