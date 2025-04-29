package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TgBotToken string
	TgAdminID  int
	Database   Database
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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

	return &Config{
		TgBotToken: tgBotToken,
		TgAdminID:  adminID,
		Database: Database{
			Host:     "postgres",
			Port:     "5432",
			Username: "novich",
			Password: "novich",
			DBName:   "vpndb",
			SSLMode:  "disable",
		},
	}
}
