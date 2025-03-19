package service

import (
	"bytes"
	"context"
	"fmt"
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

func (s *VPNService) Create(ctx context.Context, username string, chatID int) (c string, err error) {
	defer func() { err = e.WrapIfErr("can't create text config: %v", err) }()
	publicUserKey, privateUserKey, err := generateKey()
	if err != nil {
		return "", err
	}

	// ip, err := s.getNextIP(ctx)
	// if err != nil {
	// 	return "", err
	// }
	ip := "10.0.0.4"
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
		IP:         ip,
		IsActive:   true,
	}); err != nil {
		return "", err
	}

	config, err := s.conf.GenerateConfig(privateUserKey, publicUserKey, ip)
	if err != nil {
		return "", err
	}

	return config, nil
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
