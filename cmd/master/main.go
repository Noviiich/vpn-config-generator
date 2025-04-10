package main

import (
	"fmt"
	"log"

	config "github.com/Noviiich/vpn-config-generator/configs"
	tgClient "github.com/Noviiich/vpn-config-generator/internal/bot/clients/telegram"
	event_consumer "github.com/Noviiich/vpn-config-generator/internal/bot/consumer/event-consumer"
	"github.com/Noviiich/vpn-config-generator/internal/bot/events/telegram"
	"github.com/Noviiich/vpn-config-generator/internal/master/handler"
	"github.com/Noviiich/vpn-config-generator/internal/master/service"
	"github.com/Noviiich/vpn-config-generator/internal/master/storage/postgres"
	"github.com/labstack/echo/v4"
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

	handlers := handler.NewActionHandler(vpnService)
	e := echo.New()
	e.GET("/:since", handlers.GetActionsHandler)

	consumerErrChan := make(chan error)
	handlerErrChan := make(chan error)

	go func() {
		consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
		if err := consumer.Start(); err != nil {
			consumerErrChan <- fmt.Errorf("event consumer stopped: %v", err)
		}
	}()

	go func() {
		if err := e.Start(":8080"); err != nil {
			handlerErrChan <- fmt.Errorf("handler stopped: %v", err)
		}
	}()

	select {
	case err := <-consumerErrChan:
		log.Fatal(err)
	case err := <-handlerErrChan:
		log.Fatal(err)
	}
}
