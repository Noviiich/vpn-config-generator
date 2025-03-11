package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Noviiich/vpn-config-generator/config"
	"github.com/Noviiich/vpn-config-generator/internal/adapters/api"
)

func main() {
	cfg := config.Load()
	client := api.NewTimeWebClient("/api/v1", "api.timeweb.cloud", cfg.ServerToken)
	bodyText, err := client.ServerInfo(context.Background(), "4383899")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

}
