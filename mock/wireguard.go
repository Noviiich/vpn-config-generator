package mock

import (
	"fmt"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

type WGMockManager struct {
	configPath       string
	PublicServerKey  string
	PrivateServerKey string
	IPAddrServer     string
}

func NewWGMockManager(configPath string) *WGMockManager {
	return &WGMockManager{
		configPath:       configPath,
		PublicServerKey:  "public",
		PrivateServerKey: "private",
		IPAddrServer:     "ipAddr",
	}
}

func (wg *WGMockManager) GetConfig(privateUserKey string, ipAddrUser string) (string, error) {
	return wg.createUserConfig(privateUserKey, ipAddrUser), nil
}

func (wg *WGMockManager) GenerateConfig(privateUserKey string, publicUserKey string, ipAddrUser string) (string, error) {
	err := wg.changeBaseConfig(publicUserKey, ipAddrUser)
	if err != nil {
		return "", fmt.Errorf("can't change base config: %v", err)
	}

	return wg.createUserConfig(privateUserKey, ipAddrUser), nil
}

func (wg *WGMockManager) createUserConfig(privateUserKey string, ipAddrUser string) string {
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

func (wg *WGMockManager) changeBaseConfig(userPublicKey, ipAddrUser string) (err error) {
	defer func() { err = e.WrapIfErr("can't change base config", err) }()
	// 	configEntry := fmt.Sprintf(`[Peer]
	// PublicKey = %s
	// AllowedIPs = %s\n`, userPublicKey, ipAddrUser)

	return nil
}
