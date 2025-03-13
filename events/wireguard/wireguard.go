package service

import (
	"bytes"
	"fmt"
	"os/exec"
)

type WGManager struct {
	configPath string
}

func NewWGManager(configPath string) *WGManager {
	return &WGManager{configPath: configPath}
}

func (g *WGManager) GenerateConfig(privateKey, ip string) string {
	publicServerByte, err := exec.Command("cat", "/etc/wireguard/server_publickey").Output()
	if err != nil {
		fmt.Printf("Ошибка: %v\nВывод: %s", err, publicServerByte)
	}
	public := string(bytes.TrimSpace(publicServerByte))

	return fmt.Sprintf(`
		[Interface]
		PrivateKey = %s
		Address = %s/24
		DNS = 8.8.8.8

		[Peer]
		PublicKey = %s
		Endpoint = %s:51820
		AllowedIPs = 0.0.0.0/0`, privateKey, ip, public, "194.87.27.77")
}

func (wg *WGManager) AddClient(publicKey, ip string) error {
	configEntry := fmt.Sprintf("\n[Peer]\nPublicKey = %s\nAllowedIPs = %s/32\n", publicKey, ip)

	// Добавление в файл конфигурации
	err := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' >> %s", configEntry, wg.configPath)).Run()
	if err != nil {
		return err
	}

	// Перезагрузка WireGuard
	err = exec.Command("sudo", "wg-quick", "down", "wg0").Run()
	if err != nil {
		return err
	}

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
