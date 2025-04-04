package main

import (
	"log"

	config "github.com/Noviiich/vpn-config-generator/configs"
	tgClient "github.com/Noviiich/vpn-config-generator/internal/bot/clients/telegram"
	event_consumer "github.com/Noviiich/vpn-config-generator/internal/bot/consumer/event-consumer"
	"github.com/Noviiich/vpn-config-generator/internal/bot/events/telegram"
	"github.com/Noviiich/vpn-config-generator/internal/master/service"
	"github.com/Noviiich/vpn-config-generator/internal/master/storage/postgres"
	_ "github.com/lib/pq"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	cfg := config.Load()
	repo := postgres.New("novich", "novich", "vpndb")
	vpnService := service.NewVPNService(repo)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken),
		vpnService,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatalf("неисправная ошибка: %v", err)
	}
}
