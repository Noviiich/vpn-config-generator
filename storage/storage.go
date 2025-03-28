package storage

import (
	"context"
	"time"
)

type Storage interface {
	Users
	Devices
	SubscriptionTypes
}

type Users interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsers(ctx context.Context) ([]User, error)
	IsExistsUser(ctx context.Context, id int) (bool, error)
}

type Devices interface {
	CreateDevice(ctx context.Context, device *Device) error
	GetDevices(ctx context.Context, userID int) ([]Device, error)
	IsExistsDevice(ctx context.Context, chatID int) (bool, error)
}

type IPPool interface {
}

type SubscriptionTypes interface {
	ActivateSubscription(ctx context.Context, userID, typeID int) error
	DeactivateSubscription(ctx context.Context, userID int) error
	CheckSubscription(ctx context.Context, userID int) (*Subscription, error)
}

type Subscriptions interface {
}

type User struct {
	ID         int       `db:"id"`
	TelegramID int       `db:"telegram_id"`
	Username   string    `db:"username"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Device struct {
	ID         string `db:"id"`
	UserID     int    `db:"user_id"`
	PrivateKey string `db:"private_key"`
	PublicKey  string `db:"public_key"`
	CreatedAt  time.Time
	LastActive time.Time
	IsActive   bool `db:"is_active"`
}

type IPAddress struct {
	ID          int
	DeviceID    int
	IP          string
	IsAvailable bool
}

type SubscriptionType struct {
	ID         int
	Name       string
	Duration   time.Duration
	MaxDevices int
}

type Subscription struct {
	UserID     int `db:"user_id"`
	TypeID     int
	IsActive   bool `db:"is_active"`
	StartDate  time.Time
	ExpiryDate time.Time `db:"expiry_date"`
}
