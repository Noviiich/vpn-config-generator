package service

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/Noviiich/vpn-config-generator/storage"
	"github.com/Noviiich/vpn-config-generator/vpnconfig"
)

type VPNService struct {
	conf vpnconfig.VPNConfig
	repo storage.Storage
}

func NewVPNService(conf vpnconfig.VPNConfig, repo storage.Storage) *VPNService {
	return &VPNService{
		conf: conf,
		repo: repo,
	}
}

func (s *VPNService) Create(ctx context.Context, username string, chatID int) (string, error) {
	publicUserKey, privateUserKey, err := generateKey()
	if err != nil {
		return "", err
	}

	exists, err := s.repo.IsExistsUser(chatID)
	if err != nil {
		return "", err
	}

	if !exists {
		if err := s.repo.CreateUser(storage.User{
			TelegramID:         chatID,
			Username:           username,
			SubscriptionActive: true,
			SubscriptionExpiry: time.Time{},
		}); err != nil {
			return "", err
		}
	}

	if err := s.repo.CreateDevice(&storage.Device{
		UserID:     chatID,
		PrivateKey: privateUserKey,
		PublicKey:  publicUserKey,
		IP:         "10.0.0.3",
		IsActive:   true,
	}); err != nil {
		return "", err
	}

	config, err := s.conf.GenerateConfig(privateUserKey, publicUserKey, "10.0.0.3")
	if err != nil {
		return "", err
	}

	return config, nil
}

func generateKey() (private string, public string, err error) {
	privateKeyBytes, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		return "", "", err
	}
	private = string(bytes.TrimSpace(privateKeyBytes))

	cmd := exec.Command("wg", "pubkey")
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(private))
	}()

	publicKeyBytes, err := cmd.Output()
	public = string(publicKeyBytes[:len(publicKeyBytes)-1])

	return
}
