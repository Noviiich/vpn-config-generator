package service

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

type WGManager struct {
	configPath       string
	PublicServerKey  string
	PrivateServerKey string
	IPAddrServer     string
}

func NewWGManager(configPath string) *WGManager {
	publicServerByte, err := exec.Command("cat", "/etc/wireguard/server_publickey").Output()
	if err != nil {
		log.Fatalf("can't get public server key: %v", err)
	}
	public := string(bytes.TrimSpace(publicServerByte))

	privateServerByte, err := exec.Command("cat", "/etc/wireguard/server_privatekey").Output()
	if err != nil {
		log.Fatalf("can't get private server key: %v", err)
	}
	private := string(bytes.TrimSpace(privateServerByte))

	IPAdressServerByte, err := exec.Command("ip", "-o", "route", "get", "to", "8.8.8.8").Output()
	if err != nil {
		log.Fatalf("can't get ip address server: %v", err)
	}
	ipAddr := string(bytes.TrimSpace(IPAdressServerByte))

	return &WGManager{
		configPath:       configPath,
		PublicServerKey:  public,
		PrivateServerKey: private,
		IPAddrServer:     ipAddr,
	}
}

func (wg *WGManager) GenerateConfig(privateUserKey string, publicUserKey string, ipAddrUser string) (string, error) {
	err := wg.changeBaseConfig(publicUserKey, ipAddrUser)
	if err != nil {
		return "", fmt.Errorf("can't change base config: %v", err)
	}

	return wg.createUserConfig(privateUserKey, ipAddrUser), nil
}

func (wg *WGManager) createUserConfig(privateUserKey string, ipAddrUser string) string {
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/24
DNS = 8.8.8.8

[Peer]
PublicKey = %s
Endpoint = %s:51820
AllowedIPs = 0.0.0.0/0`, privateUserKey, ipAddrUser, wg.PublicServerKey, wg.IPAddrServer)
	return config
}

func (wg *WGManager) changeBaseConfig(userPublicKey, ipAddrUser string) (err error) {
	defer func() { err = e.WrapIfErr("can't change base config", err) }()
	configEntry := fmt.Sprintf(`[Peer]
PublicKey = %s
AllowedIPs = %s/32`, userPublicKey, ipAddrUser)

	// Добавление в файл конфигурации
	err = exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' >> %s", configEntry, wg.configPath)).Run()
	if err != nil {
		return err
	}

	// Перезагрузка WireGuard
	err = exec.Command("sudo", "wg-quick", "down", "wg0").Run()
	if err != nil {
		return err
	}

	// Запуск WireGuard
	err = exec.Command("sudo", "wg-quick", "up", "wg0").Run()
	if err != nil {
		return err
	}

	return nil
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
