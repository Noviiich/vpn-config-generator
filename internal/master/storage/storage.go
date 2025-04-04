package storage

import (
	"context"
	"time"
)

type Storage interface {
	Users
}

type Users interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsers(ctx context.Context) ([]User, error)
}

type User struct {
	ID         int
	TelegramID int
	Username   string
}

type SubscriptionType struct {
	ID         int           `db:"id"`
	Name       string        `db:"name"`
	Duration   time.Duration `db:"duration"`
	MaxDevices int           `db:"max_devices"`
}

type Subscription struct {
	UserID     int       `db:"user_id"`
	TypeID     int       `db:"type_id"`
	IsActive   bool      `db:"is_active"`
	StartDate  time.Time `db:"start_date"`
	ExpiryDate time.Time `db:"expiry_date"`
}

type ActionType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Action struct {
	ID        int
	ActionID  int
	UserID    int
	CreatedAt time.Time
}
