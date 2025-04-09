package postgres

import (
	"context"
	"database/sql"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

func (s *Storage) CreateAction(ctx context.Context, action *storage.Action) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
			INSERT INTO actions (action_id, user_id) 
			VALUES ($1, $2)
			RETURNING id`

	err = tx.QueryRowContext(ctx, query, action.ActionID, action.UserID).Scan(&action.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetActions(ctx context.Context, telegramID int) ([]storage.Action, error) {
	query := `
			SELECT *
			FROM actions`

	var actions []storage.Action
	err := s.db.SelectContext(ctx, &actions, query)
	if err != nil {
		return nil, err
	}

	if actions == nil {
		return nil, sql.ErrNoRows
	}

	return actions, nil
}
