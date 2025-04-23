package main

import (
	"github.com/Noviiich/vpn-config-generator/internal/config"
	"github.com/Noviiich/vpn-config-generator/internal/slave/handlers"
	"github.com/Noviiich/vpn-config-generator/internal/slave/protocol/wireguard"
	"github.com/Noviiich/vpn-config-generator/internal/slave/service"
	"github.com/Noviiich/vpn-config-generator/internal/slave/storage/postgres"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	repo := postgres.New(
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	vpnConfig := wireguard.NewWGManager("/etc/wireguard/wg0.conf")
	vpnService := service.NewVPNService(vpnConfig, repo)
	handler := handlers.New(vpnService)

	e := echo.New()
	e.GET("/config", handler.GetVPNService)
	e.Logger.Fatal(e.Start(":8080"))

}
