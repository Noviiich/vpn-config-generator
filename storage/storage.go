package storage

import (
	"context"
	"time"
)

type Storage interface {
	InitDB()
	CreateDevice(device *Device) error
	CreateUser(user User) error
	IsExistsUser(telegramID int) (bool, error)
	GetIP(ctx context.Context) (string, error)
	UpdateIP(ctx context.Context, newIP string) error
}

type User struct {
	TelegramID         int
	Username           string
	Devices            []Device
	SubscriptionActive bool
	SubscriptionExpiry time.Time
}

type Device struct {
	ID         string
	UserID     int
	PrivateKey string
	PublicKey  string
	IP         string
	IsActive   bool
}
