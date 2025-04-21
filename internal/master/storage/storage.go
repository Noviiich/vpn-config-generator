package storage

import (
	"context"
	"time"
)

type Storage interface {
	Users
	Subscriptions
	Actions
}

type Users interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
	GetUsers(ctx context.Context) ([]User, error)
}

type Subscriptions interface {
	CreateSubscription(ctx context.Context, userID, typeID int) error
	DeleteSubscription(ctx context.Context, userID int) error
	GetSubscription(ctx context.Context, userID int) (*Subscription, error)
	GetSubscriptions(ctx context.Context, userID int) ([]Subscription, error)
}

type Actions interface {
	CreateAction(ctx context.Context, action *Action) error
	GetActions(ctx context.Context) ([]Action, error)
	GetActionsByTime(ctx context.Context, since time.Time) ([]Action, error)
}

type User struct {
	ID         int       `db:"id"`
	TelegramID int       `db:"telegram_id"`
	Username   string    `db:"username"`
	CreatedAt  time.Time `db:"created_at"`
}

type SubscriptionType struct {
	ID         int           `db:"id"`
	Name       string        `db:"name"`
	Duration   time.Duration `db:"duration"`
	MaxDevices int           `db:"max_devices"`
}

type Subscription struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	SubscriptionID int       `db:"subscription_id"`
	ExpiryDate     time.Time `db:"expiry_date"`
}

type ActionType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Action struct {
	ID        int       `db:"id"`
	ActionID  int       `db:"action_id"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
