package storage

import (
	"context"
	"time"
)

type Storage interface {
	InitDB()
	CreateDevice(ctx context.Context, device *Device) error
	CreateUser(ctx context.Context, user *User) error
	IsExistsUser(ctx context.Context, telegramID int) (bool, error)
	GetIP(ctx context.Context) (string, error)
	UpdateIP(ctx context.Context, newIP string) error
	GetDevice(ctx context.Context, telegramID int) (Device, error)
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
