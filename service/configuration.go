package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *VPNService) createConfig(device *storage.Device) (c string, err error) {
	configText, err := s.conf.GenerateConfig(device.PrivateKey, device.PublicKey, device.IP)
	if err != nil {
		return "", e.Wrap("can't create config", err)
	}
	return configText, nil
}

func (s *VPNService) GetConfig(ctx context.Context, chatID int, username string) (con string, err error) {
	defer func() { err = e.WrapIfErr("can't get config", err) }()
	user, err := s.GetUser(ctx, chatID, username)
	if err != nil {
		return "", err
	}

	device, err := s.GetDevice(ctx, user.TelegramID)
	if err != nil {
		return "", err
	}

	conf, err := s.createConfig(device)
	if err != nil {
		return "", err
	}

	return conf, nil
}
