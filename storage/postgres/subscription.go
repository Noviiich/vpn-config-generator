package postgres

import (
	"context"
	"fmt"
)

func (r *Storage) DeactivateSubscription(ctx context.Context, userID int) error {
	query := `
		UPDATE users
		SET subscription_active = FALSE
		WHERE telegram_id = $1;
	`
	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deactivating subscription for user %d: %w", userID, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription deactivated: user with id %d not found", userID)
	}

	return nil
}
