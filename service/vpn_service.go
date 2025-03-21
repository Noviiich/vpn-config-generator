package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
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
		msg := fmt.Sprintf(`Вы молодец, у вас есть подписка!!!
		Ваша подписка истекает через %s`, remaining.Truncate(time.Second))
		return msg, nil
	}

	return "У вас не подписки. Не расстраивайтесь, вы все еще можете ее оформить", nil

}

func (s *VPNService) isExistsUser(ctx context.Context, chatID int) (ex bool, err error) {
	defer func() { err = e.WrapIfErr("can't is exists user", err) }()
	exists, err := s.repo.IsExistsUser(ctx, chatID)
	if err != nil {
		return exists, err
	}
	return exists, nil
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

func (s *VPNService) createNewUser(ctx context.Context, username string, chatID int) (u *storage.User, err error) {
	defer func() { err = e.WrapIfErr("can't create new user", err) }()
	user := &storage.User{
		TelegramID:         chatID,
		Username:           username,
		SubscriptionActive: true,
		SubscriptionExpiry: time.Now().AddDate(0, 1, 0),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *VPNService) createNewDevice(ctx context.Context, userID int) (d *storage.Device, err error) {
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

func (s *VPNService) getNextIP(ctx context.Context) (ip string, err error) {
	defer func() { err = e.WrapIfErr("can't get next ip: %v", err) }()
	lastIP, err := s.repo.GetIP(ctx)
	if err != nil {
		return "", err
	}

	ipFormat := net.ParseIP(lastIP)
	if ipFormat == nil {
		return "", err
	}

	ipParts := ipFormat.To4()
	if ipParts == nil {
		return "", err
	}

	lastOctet, _ := strconv.Atoi(fmt.Sprintf("%d", ipParts[3]))
	if lastOctet >= 254 { // Максимум x.x.x.254
		return "", err
	}

	newIP := fmt.Sprintf("%d.%d.%d.%d",
		ipParts[0],
		ipParts[1],
		ipParts[2],
		lastOctet+1,
	)

	if err = s.repo.UpdateIP(ctx, newIP); err != nil {
		return "", err
	}

	return newIP, nil
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
