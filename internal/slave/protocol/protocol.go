package protocol

type Protocol interface {
	GenerateConfig(privateUserKey string, publicUserKey string, ipAddrUser string) (string, error)
	GetConfig(privateUserKey string, ipAddrUser string) (string, error)
}
