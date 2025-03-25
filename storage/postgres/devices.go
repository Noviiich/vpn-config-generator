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
	if err := s.db.GetContext(ctx, &device, query, telegramID); err != nil {
		return nil, err
	}

	return &device, nil
}

func (s *Storage) CreateDevice(ctx context.Context, device *storage.Device) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO devices (user_id, private_key, public_key, ip, is_active) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	row := tx.QueryRowContext(ctx, query,
		device.UserID, device.PrivateKey, device.PublicKey, device.IP, device.IsActive)
	err = row.Scan(&device.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) IsExistsDevice(ctx context.Context, telegramID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1)`
	err := s.db.QueryRowContext(ctx, query, telegramID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
