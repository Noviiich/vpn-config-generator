package service

import "github.com/Noviiich/vpn-config-generator/internal/master/storage"

type VPNService struct {
	db storage.Storage
}

func NewVPNService(db storage.Storage) *VPNService {
	return &VPNService{
		db: db,
	}
}
