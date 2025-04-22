package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Noviiich/vpn-config-generator/internal/config"
	tgClient "github.com/Noviiich/vpn-config-generator/internal/master/clients/telegram"
	event_consumer "github.com/Noviiich/vpn-config-generator/internal/master/consumer/event-consumer"
	subscription_consumer "github.com/Noviiich/vpn-config-generator/internal/master/consumer/subscription-consumer"
	"github.com/Noviiich/vpn-config-generator/internal/master/events/telegram"
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
	repo := postgres.New(
		cfg.Database.Username,
		cfg.Database.Password,
		"master",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	vpnService := service.NewVPNService(repo)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken, cfg.TgAdminID),
		vpnService,
	)

	log.Print("service started")

	consumerErrChan := make(chan error)
	subCheckErrChan := make(chan error)

	go func() {
		consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
		if err := consumer.Start(); err != nil {
			consumerErrChan <- fmt.Errorf("event consumer stopped: %v", err)
		}
	}()

	go func() {
		checkSub := subscription_consumer.New(repo, 1*time.Hour)
		if err := checkSub.Start(); err != nil {
			subCheckErrChan <- fmt.Errorf("subscription checker stopped: %v", err)
		}
	}()

	select {
	case err := <-consumerErrChan:
		log.Fatal(err)
	case err := <-subCheckErrChan:
		log.Fatal(err)
	}
}
