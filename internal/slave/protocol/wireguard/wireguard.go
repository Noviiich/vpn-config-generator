package wireguard

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/Noviiich/vpn-config-generator/internal/lib/e"
)

type Wireguard struct {
	configPath       string
	PublicServerKey  string
	PrivateServerKey string
	IPAddrServer     string
}

func NewWGManager(configPath string) *Wireguard {
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

	IPAdressServerByte, err := exec.Command("sh", "-c", "ip -o route get 8.8.8.8 | awk '/src/ {print $7}'").Output()
	if err != nil {
		log.Fatalf("can't get ip address server: %v", err)
	}
	ipAddr := string(bytes.TrimSpace(IPAdressServerByte))

	return &Wireguard{
		configPath:       configPath,
		PublicServerKey:  public,
		PrivateServerKey: private,
		IPAddrServer:     ipAddr,
	}
}

func (wg *Wireguard) GetConfig(privateUserKey string, ipAddrUser string) (string, error) {
	return wg.createUserConfig(privateUserKey, ipAddrUser), nil
}

func (wg *Wireguard) GenerateConfig(privateUserKey string, publicUserKey string, ipAddrUser string) (string, error) {
	err := wg.changeBaseConfig(publicUserKey, ipAddrUser)
	if err != nil {
		return "", fmt.Errorf("can't change base config: %v", err)
	}

	return wg.createUserConfig(privateUserKey, ipAddrUser), nil
}

func (wg *Wireguard) createUserConfig(privateUserKey string, ipAddrUser string) string {
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/32
DNS = 8.8.8.8

[Peer]
PublicKey = %s
Endpoint = %s:51820
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 20`, privateUserKey, ipAddrUser, wg.PublicServerKey, wg.IPAddrServer)
	return config
}

func (wg *Wireguard) changeBaseConfig(userPublicKey, ipAddrUser string) (err error) {
	defer func() { err = e.WrapIfErr("can't change base config", err) }()
	configEntry := fmt.Sprintf(`[Peer]
PublicKey = %s
AllowedIPs = %s\n`, userPublicKey, ipAddrUser)

	// Добавление в файл конфигурации
	err = exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' >> %s", configEntry, wg.configPath)).Run()
	if err != nil {
		return err
	}

	// остановка WireGuard
	err = exec.Command("systemctl", "stop", "wg-quick@wg0").Run()
	if err != nil {
		return err
	}

	// старт WireGuard
	err = exec.Command("systemctl", "start", "wg-quick@wg0").Run()
	if err != nil {
		return err
	}

	return nil
}
