package vpnconfig

type VPNConfig interface {
	GenerateConfig(string, string, string) (string, error)
}
