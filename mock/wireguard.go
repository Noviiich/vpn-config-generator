package mock

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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
		IPAddrServer:     "10.0.0.1",
	}
}

func (wg *WGMockManager) GetConfig(privateUserKey string, ipAddrUser string) (string, error) {
	return wg.createUserConfig(privateUserKey, ipAddrUser), nil
}

func (wg *WGMockManager) GenerateConfig(ipAddrUser string) (string, error) {
	private, _, err := generateKeys()
	if err != nil {
		return "", err
	}

	return wg.createUserConfig(private, ipAddrUser), nil
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

// func (wg *WGMockManager) changeBaseConfig(userPublicKey, ipAddrUser string) (err error) {
// 	defer func() { err = e.WrapIfErr("can't change base config", err) }()
// 	// 	configEntry := fmt.Sprintf(`[Peer]
// 	// PublicKey = %s
// 	// AllowedIPs = %s\n`, userPublicKey, ipAddrUser)

// 	return nil
// }

func generateKeys() (string, string, error) {
	key, err := wgtypes.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	private := key.String()
	public := key.PublicKey().String()

	return private, public, nil
}
