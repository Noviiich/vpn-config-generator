package postgres

import (
	"context"
	"fmt"
)

func (r *Storage) DeactivateSubscription(ctx context.Context, userID int) error {
	query := `UPDATE subscription
			  SET is_active = FALSE
		      WHERE user_id = $1;`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deactivating subscription for user %d: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for user %d: %w", userID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription deactivated: user with id %d not found", userID)
	}

	return nil
}

func (s *Storage) CreateSubscription(ctx context.Context, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `INSERT INTO subscription 
			  (telegram_id) VALUES ($1)`

	_, err = tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
