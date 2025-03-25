package storage

import (
	"context"
	"time"
)

type Storage interface {
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
	GetUsers(ctx context.Context) ([]User, error)
}

type User struct {
	TelegramID         int       `db:"telegram_id"`
	Username           string    `db:"username"`
	SubscriptionActive bool      `db:"subscription_active"`
	SubscriptionExpiry time.Time `db:"subscription_expiry"`
}

type Device struct {
	ID         string `db:"id"`
	UserID     int    `db:"user_id"`
	PrivateKey string `db:"private_key"`
	PublicKey  string `db:"public_key"`
	IP         string `db:"ip"`
	IsActive   bool   `db:"is_active"`
}

type Subscription struct {
	UserID     int       `db:"user_id"`
	IsActive   bool      `db:"is_active"`
	ExpiryDate time.Time `db:"expiry_date"`
}
