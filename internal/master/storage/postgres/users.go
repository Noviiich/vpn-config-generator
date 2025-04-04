package postgres

import (
	"context"
	"database/sql"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
			INSERT INTO users (telegram_id, username) 
			VALUES ($1, $2)
			RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		user.TelegramID, user.Username).Scan(&user.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetUser(ctx context.Context, telegramID int) (*storage.User, error) {
	query := `
			SELECT *
			FROM users
			WHERE telegram_id = $1
			LIMIT 1`

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
			SET username = $1
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

	if users == nil {
		return nil, sql.ErrNoRows
	}

	return users, nil
}
