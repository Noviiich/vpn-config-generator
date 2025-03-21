package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *VPNService) createConfig(ctx context.Context, device *storage.Device) (c string, err error) {
	configText, err := s.conf.GenerateConfig(device.PrivateKey, device.PublicKey, device.IP)
	if err != nil {
		return "", e.Wrap("can't create config", err)
	}
	return configText, nil
}

func (s *VPNService) GetConfig(ctx context.Context, chatID int, username string) (con string, err error) {
	defer func() { err = e.WrapIfErr("can't get config", err) }()
	userExists, err := s.isExistsUser(ctx, chatID)
	if err != nil {
		return "", err
	}
	if !userExists {
		err = s.CreateUser(ctx, username, chatID)
		if err != nil {
			return "", err
		}

		device, err := s.createDevice(ctx, chatID)
		if err != nil {
			return "", err
		}

		conf, err := s.createConfig(ctx, device)
		if err != nil {
			return "", err
		}

		return conf, nil
	}

	exDevice, err := s.isExistsDevice(ctx, chatID)
	if err != nil {
		return "", err
	}
	if !exDevice {
		device, err := s.createDevice(ctx, chatID)
		if err != nil {
			return "", err
		}

		conf, err := s.createConfig(ctx, device)
		if err != nil {
			return "", err
		}

		return conf, nil
	}

	device, err := s.getDevice(ctx, chatID)
	if err != nil {
		return "", err
	}

	conf, err := s.createConfig(ctx, device)
	if err != nil {
		return "", err
	}

	return conf, nil
}
