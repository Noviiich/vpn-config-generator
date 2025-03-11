package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerToken string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	serverToken := os.Getenv("VPN_SERVER_TOKEN")

	if serverToken == "" {
		log.Fatal("VPN_SERVER_TOKEN is not specified")
	}
	return &Config{ServerToken: serverToken}
}
