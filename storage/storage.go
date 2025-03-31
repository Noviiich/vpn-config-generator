package storage

import (
	"context"
	"time"
)

type Storage interface {
	Users
	Devices
	SubscriptionTypes
	LastIPs
	IPPool
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
	GetIpIsNull(ctx context.Context) (*IPAddress, error)
	CreateIP(ctx context.Context, IPAdderss *IPAddress) error
}

type SubscriptionTypes interface {
	ActivateSubscription(ctx context.Context, userID, typeID int) error
	DeactivateSubscription(ctx context.Context, userID int) error
	CheckSubscription(ctx context.Context, userID int) (*Subscription, error)
}

type Subscriptions interface {
}

type LastIPs interface {
	GetLastIP(ctx context.Context, subnet string) (*LastIP, error)
	UpdateLastIP(ctx context.Context, newIP string, subnet string) error
}

type User struct {
	ID         int       `db:"id"`
	TelegramID int       `db:"telegram_id"`
	Username   string    `db:"username"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Device struct {
	ID         string    `db:"id"`
	UserID     int       `db:"user_id"`
	TypeId     int       `db:"type_id"`
	PrivateKey string    `db:"private_key"`
	PublicKey  string    `db:"public_key"`
	CreatedAt  time.Time `db:"created_at"`
	LastActive time.Time `db:"last_active"`
	IsActive   bool      `db:"is_active"`
}

type IPAddress struct {
	ID          int    `db:"id"`
	DeviceID    int    `db:"device_id"`
	IP          string `db:"ip"`
	IsAvailable bool   `db:"is_available"`
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

type LastIP struct {
	Subnet string `db:"subnet"`
	IP     string `db:"ip"`
}
