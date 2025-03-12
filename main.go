package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Noviiich/vpn-config-generator/config"
	"github.com/Noviiich/vpn-config-generator/internal/adapters/api"
	"github.com/Noviiich/vpn-config-generator/internal/core/service"
)

func main() {
	cfg := config.Load()
	service := service.NewWGManager("/etc/wireguard/wg0.conf")
	private, public, err := service.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	err = service.AddClient(public, "10.0.0.2")
	if err != nil {
		log.Fatal(err)
	}
	config := service.GenerateConfig(private, "10.0.0.2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
	client := api.NewTimeWebClient("/api/v1", "api.timeweb.cloud", cfg.ServerToken)
	bodyText, err := client.ServerInfo(context.Background(), "4383899")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

}
