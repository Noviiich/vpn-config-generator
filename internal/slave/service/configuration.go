package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/slave/lib/e"
	"github.com/Noviiich/vpn-config-generator/internal/slave/storage"
)

func (s *VPNService) createConfig(device *storage.Device) (c string, err error) {
	configText, err := s.protocol.GenerateConfig(device.PrivateKey, device.PublicKey, device.IP)
	if err != nil {
		return "", e.Wrap("can't create config", err)
	}
	return configText, nil
}

func (s *VPNService) GetConfig(ctx context.Context, userID int) (con string, err error) {

	device, err := s.GetDevice(ctx, userID)
	if err != nil {
		return "", err
	}

	conf, err := s.createConfig(device)
	if err != nil {
		return "", err
	}

	return conf, nil
}
