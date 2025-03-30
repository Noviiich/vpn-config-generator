package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
			INSERT INTO users (telegram_id, username) 
			VALUES ($1, $2)
			RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		user.TelegramID, user.Username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetUser(ctx context.Context, telegramID int) (*storage.User, error) {
	query := `
			SELECT *
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
	query := `
			UPDATE users 
			SET username = $1, updated_at = NOW() 
			WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query,
		user.Username, user.TelegramID)
	return err
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	query := `
			DELETE FROM users WHERE id = $1`
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

	return users, nil
}

func (s *Storage) IsExistsUser(ctx context.Context, telegramID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)`
	err := s.db.QueryRowContext(ctx, query, telegramID).Scan(&exists)
	return exists, err
}
