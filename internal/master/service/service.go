package service

import (
	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

type VPNService struct {
	repo storage.Storage
}

func NewVPNService(repo storage.Storage) *VPNService {
	return &VPNService{repo: repo}
}
