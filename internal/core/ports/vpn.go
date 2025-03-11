package ports

type VPNConfigProvider interface {
	GenerateConfig(privateKey, publicKey, ip string) string
}
