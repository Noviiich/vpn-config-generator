package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *Storage) CreateDevice(ctx context.Context, device *storage.Device) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
			INSERT INTO devices (user_id, private_key, public_key) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, created_at, is_active`

	err = tx.QueryRowContext(ctx, query,
		device.UserID, device.PrivateKey, device.PublicKey).Scan(
		&device.ID,
		&device.CreatedAt,
		&device.IsActive,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetDevices(ctx context.Context, userID int) ([]storage.Device, error) {
	query := `SELECT *
              FROM devices 
              WHERE user_id = $1`

	var devices []storage.Device
	err := s.db.SelectContext(ctx, &devices, query)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (s *Storage) IsExistsDevice(ctx context.Context, chatID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1)`
	err := s.db.QueryRowContext(ctx, query, chatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
