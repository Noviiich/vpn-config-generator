package vpnconfig

type VPNConfig interface {
	GenerateConfig(privateUserKey string, publicUserKey string, ipAddrUser string) (string, error)
}
