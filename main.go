package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tgClient "github.com/Noviiich/vpn-config-generator/clients/telegram"
	"github.com/Noviiich/vpn-config-generator/config"
	event_consumer "github.com/Noviiich/vpn-config-generator/consumer/event-consumer"
	subscription_consumer "github.com/Noviiich/vpn-config-generator/consumer/subscription-consumer"
	"github.com/Noviiich/vpn-config-generator/events/telegram"
	"github.com/Noviiich/vpn-config-generator/service"
	"github.com/Noviiich/vpn-config-generator/storage/postgres"
	"github.com/Noviiich/vpn-config-generator/vpnconfig/wireguard"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	cfg := config.Load()
	// service := service.NewWGManager("/etc/wireguard/wg0.conf")
	// private, public, err := service.GenerateKey()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = service.AddClient(public, "10.0.0.2")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// config := service.GenerateConfig(private, "10.0.0.2")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(config)
	// client := api.NewTimeWebClient("/api/v1", "api.timeweb.cloud", cfg.ServerToken)
	// bodyText, err := client.ServerInfo(context.Background(), "4383899")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", bodyText)
	repo := postgres.New("novich", "novich", "vpndb")
	repo.InitDB(context.Background())
	vpnConfig := wireguard.NewWGManager("/etc/wireguard/wg0.conf")
	vpnService := service.NewVPNService(vpnConfig, repo)

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
		checkSub := subscription_consumer.New(repo, 1*time.Second)
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
