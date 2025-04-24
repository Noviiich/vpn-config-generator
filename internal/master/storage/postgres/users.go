package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (telegram_id, username, subscription_active, subscription_expiry) 
		 VALUES ($1, $2, $3, $4) ON CONFLICT (telegram_id) DO NOTHING`

	_, err = tx.ExecContext(ctx, query,
		user.TelegramID,
		user.Username,
		user.SubscriptionActive,
		user.SubscriptionExpiry)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetUser(ctx context.Context, telegramID int) (*storage.User, error) {
	query := `SELECT telegram_id, username, subscription_active, subscription_expiry
              FROM users
              WHERE telegram_id = $1`

	var user storage.User
	err := s.db.GetContext(ctx, &user, query, telegramID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *storage.User) error {
	query := `UPDATE users 
			SET subscription_active = $1, subscription_expiry = $2 
			WHERE telegram_id = $3`
	_, err := s.db.ExecContext(ctx, query,
		user.SubscriptionActive, user.SubscriptionExpiry, user.TelegramID)
	return err
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE telegram_id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *Storage) GetUsers(ctx context.Context) ([]storage.User, error) {
	query := `SELECT * FROM users`

	var users []storage.User
	err := s.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	if users == nil {
		return nil, nil
	}

	return users, nil
}

func (s *Storage) IsExistsUser(ctx context.Context, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)`
	err := s.db.GetContext(ctx, &exists, query, userID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
