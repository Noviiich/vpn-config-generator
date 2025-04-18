package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// ServerToken string
	TgBotToken string
	TgAdminID  int
	// MongoToken  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// serverToken := os.Getenv("VPN_SERVER_TOKEN")

	// if serverToken == "" {
	// 	log.Fatal("VPN_SERVER_TOKEN is not specified")
	// }

	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	adminIDStr := os.Getenv("TG_ADMIN_ID")
	if adminIDStr == "" {
		log.Fatal("TG_ADMIN_ID is not specified")
	}

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		log.Fatalf("Invalid TG_ADMIN_ID: %v", err)
	}

	if tgBotToken == "" {
		log.Fatal("TG_BOT_TOKEN is not specified")
	}

	// mongoToken := os.Getenv("MONGO_CONNECTION_STRING")

	// if mongoToken == "" {
	// 	log.Fatal("MONGO_CONNECTION_STRING is not specified")
	// }
	return &Config{
		// ServerToken: serverToken,
		TgBotToken: tgBotToken,
		TgAdminID:  adminID,
		// MongoToken:  mongoToken,
	}
}
