package service

import (
	"github.com/Noviiich/vpn-config-generator/internal/core/ports"
)

type VPNService struct {
	vpnProvider ports.VPNConfigProvider
}

func NewVPNService(vpnProvider ports.VPNConfigProvider) *VPNService {
	return &VPNService{vpnProvider: vpnProvider}
}

func (s *VPNService) GenerateConfig(privateKey, publicKey, ip string) string {

}
