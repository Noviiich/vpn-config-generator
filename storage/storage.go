package storage

type Storage interface {
	Create(Config, string) (string, error)
}

type Config struct {
	ID         string
	PrivateKey string
	PublicKey  string
	IPAddress  string
}
