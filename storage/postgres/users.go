package postgres

import (
	"context"
	"fmt"

	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO users (telegram_id, username, subscription_active, subscription_expiry) 
		 VALUES ($1, $2, $3, $4) ON CONFLICT (telegram_id) DO NOTHING`,
		user.TelegramID, user.Username, user.SubscriptionActive, user.SubscriptionExpiry)
	return err
}

func (s *Storage) IsExistsUser(ctx context.Context, telegramID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)`
	err := s.pool.QueryRow(ctx, query, telegramID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("can't check if user exists: %v", err)
	}
	return exists, nil
}

func (s *Storage) GetUser(ctx context.Context, telegramID int) (*storage.User, error) {
	query := `SELECT telegram_id, username, subscription_active, subscription_expiry
              FROM users
              WHERE telegram_id = $1`
	var user storage.User
	err := s.pool.QueryRow(ctx, query, telegramID).Scan(
		&user.TelegramID,
		&user.Username,
		&user.SubscriptionActive,
		&user.SubscriptionExpiry,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *storage.User) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET subscription_active = $1, subscription_expiry = $2 WHERE telegram_id = $3`,
		user.SubscriptionActive, user.SubscriptionExpiry, user.TelegramID)
	return err
}
