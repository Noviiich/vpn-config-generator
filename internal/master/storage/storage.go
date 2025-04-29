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

	GetServers(ctx context.Context) ([]Server, error)
	AddServer(ctx context.Context, server *Server) error
	DeleteServer(ctx context.Context, serverID int) error
	GetServer(ctx context.Context, serverID int) (*Server, error)
}

type User struct {
	TelegramID         int       `db:"telegram_id" json:"telegram_id"`
	Username           string    `db:"username" json:"username"`
	SubscriptionActive bool      `db:"subscription_active" json:"subscription_active"`
	SubscriptionExpiry time.Time `db:"subscription_expiry" json:"subscription_expiry"`
}

type Server struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	CountryFlag string `db:"country_flag" json:"country_flag"`
	IPAddress   string `db:"ip_address" json:"ip_address"`
	Port        int    `db:"port" json:"port"`
	Protocol    string `db:"protocol" json:"protocol"`
}
