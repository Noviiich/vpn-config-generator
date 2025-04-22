package main

import (
	"github.com/Noviiich/vpn-config-generator/internal/config"
	"github.com/Noviiich/vpn-config-generator/internal/slave/protocol/wireguard"
	"github.com/Noviiich/vpn-config-generator/internal/slave/service"
	"github.com/Noviiich/vpn-config-generator/internal/slave/storage/postgres"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	repo := postgres.New(
		cfg.Database.Username,
		cfg.Database.Password,
		"slave",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	vpnConfig := wireguard.NewWGManager("/etc/wireguard/wg0.conf")
	_ = service.NewVPNService(vpnConfig, repo)

}
