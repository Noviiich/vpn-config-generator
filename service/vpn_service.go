package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/Noviiich/vpn-config-generator/lib/e"
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

func (s *VPNService) StatusSubscribtion(ctx context.Context, username string, chatID int) (st string, err error) {
	defer func() { err = e.WrapIfErr("can't get status subscription", err) }()

	exists, err := s.isExistsUser(ctx, chatID)
	if err != nil {
		return "", nil
	}

	if !exists {
		_, err = s.createNewUser(ctx, username, chatID)
		if err != nil {
			return "", err
		}
	}

	user, err := s.repo.GetUser(ctx, chatID)
	if err != nil {
		return "", err
	}

	if user.SubscriptionActive {
		remaining := time.Until(user.SubscriptionExpiry)
		days := int(remaining.Hours()) / 24
		hours := int(remaining.Hours()) % 24

		msg := fmt.Sprintf(`Вы молодец, у вас есть подписка!!!
Ваша подписка истекает через %d дней, %d часов`, days, hours)
		return msg, nil
	}

	return "У вас не подписки. Не расстраивайтесь, вы все еще можете ее оформить", nil

}

func (s *VPNService) Create(ctx context.Context, username string, chatID int) (c string, err error) {
	defer func() { err = e.WrapIfErr("can't create text config", err) }()
	exists, err := s.isExistsUser(ctx, chatID)
	if err != nil {
		return "", nil
	}

	if !exists {
		return s.createNewConfig(ctx, username, chatID)
	}

	return s.getConfig(ctx, chatID)
}

func (s *VPNService) getConfig(ctx context.Context, chatID int) (string, error) {
	device, err := s.repo.GetDevice(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("can't get config: %v", err)
	}

	conf, err := s.conf.GetConfig(device.PrivateKey, device.IP)
	if err != nil {
		return "", fmt.Errorf("can't get config: %v", err)
	}

	return conf, nil
}

func (s *VPNService) createNewConfig(ctx context.Context, username string, chatID int) (c string, err error) {
	defer func() { err = e.WrapIfErr("can't create new config", err) }()
	user, err := s.createNewUser(ctx, username, chatID)
	if err != nil {
		return "", err
	}

	device, err := s.createNewDevice(ctx, user.TelegramID)
	if err != nil {
		return "", err
	}

	config, err := s.conf.GenerateConfig(device.PrivateKey, device.PublicKey, device.IP)
	if err != nil {
		return "", err
	}
	return config, nil
}

func generateKey() (private string, public string, err error) {
	defer func() { err = e.WrapIfErr("can't generate keys: %v", err) }()
	cmd := exec.Command("sh", "-c", "wg genkey | tee /etc/wireguard/user1_privatekey | wg pubkey | tee /etc/wireguard/user1_publickey")
	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	publicUserByte, err := exec.Command("cat", "/etc/wireguard/user1_publickey").Output()
	if err != nil {
		log.Fatalf("can't get public user key: %v", err)
	}
	pub := string(bytes.TrimSpace(publicUserByte))

	privateUserByte, err := exec.Command("cat", "/etc/wireguard/user1_privatekey").Output()
	if err != nil {
		log.Fatalf("can't get private user key: %v", err)
	}
	pr := string(bytes.TrimSpace(privateUserByte))

	return pr, pub, nil
}
