package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/lib/e"
	"github.com/Noviiich/vpn-config-generator/internal/slave/storage"
)

func (s *VPNService) СreateDevice(ctx context.Context, userID int) (err error) {
	defer func() { err = e.WrapIfErr("can't create device", err) }()
	privateUserKey, publicUserKey, err := generateKey()
	if err != nil {
		return err
	}

	ipUser, err := s.getNextIP(ctx)
	if err != nil {
		return err
	}

	device := &storage.Device{
		UserID:     userID,
		PrivateKey: privateUserKey,
		PublicKey:  publicUserKey,
		IP:         ipUser,
		IsActive:   true,
	}
	if err := s.repo.CreateDevice(ctx, device); err != nil {
		return err
	}

	return nil
}

func (s *VPNService) GetDevice(ctx context.Context, chatID int) (d *storage.Device, err error) {
	defer func() { err = e.WrapIfErr("can't get device", err) }()

	exists, err := s.isExistsDevice(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = s.СreateDevice(ctx, chatID)
		if err != nil {
			return nil, err
		}
	}

	device, err := s.repo.GetDevice(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *VPNService) isExistsDevice(ctx context.Context, chatID int) (bool, error) {
	exists, err := s.repo.IsExistsDevice(ctx, chatID)
	if err != nil {
		return exists, e.Wrap("can't check if device exists", err)
	}
	return exists, nil
}
