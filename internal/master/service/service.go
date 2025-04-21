package service

import (
	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
	"github.com/Noviiich/vpn-config-generator/internal/master/vpn"
)

type VPNService struct {
	db  storage.Storage
	vpn vpn.VPN
}

func NewVPNService(db storage.Storage, vpn vpn.VPN) *VPNService {
	return &VPNService{
		db:  db,
		vpn: vpn,
	}
}
