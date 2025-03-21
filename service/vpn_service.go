package service

import (
	"bytes"
	"log"
	"os/exec"

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
