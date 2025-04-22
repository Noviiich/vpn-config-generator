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
	TelegramID         int
	Username           string
	SubscriptionActive bool
	SubscriptionExpiry time.Time
}
