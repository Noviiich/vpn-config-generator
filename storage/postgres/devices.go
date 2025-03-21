package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *Storage) GetDevice(ctx context.Context, telegramID int) (*storage.Device, error) {
	query := `SELECT id, user_id, private_key, public_key, ip, is_active 
              FROM devices 
              WHERE user_id = $1`

	var device storage.Device
	err := s.pool.QueryRow(ctx, query, telegramID).Scan(
		&device.ID,
		&device.UserID,
		&device.PrivateKey,
		&device.PublicKey,
		&device.IP,
		&device.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (s *Storage) CreateDevice(ctx context.Context, device *storage.Device) error {
	query := `INSERT INTO devices (user_id, private_key, public_key, ip, is_active) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := s.pool.QueryRow(ctx, query,
		device.UserID, device.PrivateKey, device.PublicKey, device.IP, device.IsActive).Scan(&device.ID)

	return err
}
