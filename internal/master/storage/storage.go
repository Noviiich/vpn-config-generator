package storage

import (
	"context"
	"time"
)

type Storage interface {
	GetUser(ctx context.Context, telegramID int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, telegramID int) error
	IsExistsUser(ctx context.Context, telegramID int) (bool, error)
	GetUsers(ctx context.Context) ([]User, error)
}

type User struct {
	TelegramID         int       `db:"telegram_id" json:"telegram_id"`
	Username           string    `db:"username" json:"username"`
	SubscriptionActive bool      `db:"subscription_active" json:"subscription_active"`
	SubscriptionExpiry time.Time `db:"subscription_expiry" json:"subscription_expiry"`
}
