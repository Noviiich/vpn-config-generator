package storage

import (
	"context"
)

type Storage interface {
	CreateDevice(ctx context.Context, device *Device) error
	GetIP(ctx context.Context) (string, error)
	UpdateIP(ctx context.Context, newIP string) error
	GetDevice(ctx context.Context, telegramID int) (*Device, error)
	IsExistsDevice(ctx context.Context, telegramID int) (bool, error)
}
type Device struct {
	ID         string `db:"id" json:"id"`
	UserID     int    `db:"user_id" json:"user_id"`
	PrivateKey string `db:"private_key" json:"private_key"`
	PublicKey  string `db:"public_key" json:"public_key"`
	IP         string `db:"ip" json:"ip"`
	IsActive   bool   `db:"is_active" json:"is_active"`
}
