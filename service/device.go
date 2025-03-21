package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *VPNService) createDevice(ctx context.Context, userID int) (d *storage.Device, err error) {
	defer func() { err = e.WrapIfErr("can't create new device", err) }()
	privateUserKey, publicUserKey, err := generateKey()
	if err != nil {
		return nil, err
	}

	ipUser, err := s.getNextIP(ctx)
	if err != nil {
		return nil, err
	}

	device := &storage.Device{
		UserID:     userID,
		PrivateKey: privateUserKey,
		PublicKey:  publicUserKey,
		IP:         ipUser,
		IsActive:   true,
	}
	if err := s.repo.CreateDevice(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}
