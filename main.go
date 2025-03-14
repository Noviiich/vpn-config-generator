package main

import (
	"log"

	tgClient "github.com/Noviiich/vpn-config-generator/clients/telegram"
	"github.com/Noviiich/vpn-config-generator/config"
	event_consumer "github.com/Noviiich/vpn-config-generator/consumer/event-consumer"
	"github.com/Noviiich/vpn-config-generator/events/telegram"
	"github.com/Noviiich/vpn-config-generator/service"
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
	vpnConfig := wireguard.NewWGManager("/etc/wireguard/wg0.conf")
	vpnService := service.NewVPNService(vpnConfig)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken),
		vpnService,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
