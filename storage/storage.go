package storage

import (
	"context"
	"time"
)

type Storage interface {
	InitDB(ctx context.Context)
	CreateDevice(ctx context.Context, device *Device) error
	GetUser(ctx context.Context, telegramID int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, telegramID int) error
	IsExistsUser(ctx context.Context, telegramID int) (bool, error)
	GetIP(ctx context.Context) (string, error)
	UpdateIP(ctx context.Context, newIP string) error
	GetDevice(ctx context.Context, telegramID int) (*Device, error)
	IsExistsDevice(ctx context.Context, telegramID int) (bool, error)
	GetUsers(ctx context.Context, telegramID int) ([]string, error)
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
