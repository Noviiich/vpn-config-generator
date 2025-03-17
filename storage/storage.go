package storage

import "time"

type Storage interface {
}

type User struct {
	TelegramID         uint
	Username           string
	Devices            []Device
	SubscriptionActive bool
	SubscriptionExpiry time.Time
}

type Device struct {
	ID         string
	UserID     uint
	PrivateKey string
	PublicKey  string
	IP         string
	IsActive   bool
}
