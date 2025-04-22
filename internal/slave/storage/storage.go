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
	ID         string
	UserID     int
	PrivateKey string
	PublicKey  string
	IP         string
	IsActive   bool
}
