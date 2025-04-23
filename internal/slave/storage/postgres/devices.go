package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/slave/storage"
)

func (s *Storage) CreateDevice(ctx context.Context, user *storage.Device) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `INSERT INTO devices (user_id, private_key, public_key, ip, is_active) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		user.UserID, user.PrivateKey, user.PublicKey, user.IP, user.IsActive).Scan(&user.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetDevice(ctx context.Context, userID int) (*storage.Device, error) {
	query := `SELECT *
              FROM devices 
              WHERE user_id = $1`

	var device storage.Device
	err := s.db.GetContext(ctx, &device, query, userID)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (s *Storage) IsExistsDevice(ctx context.Context, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM devices WHERE user_id = $1)`
	err := s.db.GetContext(ctx, &exists, query, userID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
