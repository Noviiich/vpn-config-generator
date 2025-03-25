package postgres

import (
	"context"
	"fmt"
)

func (r *Storage) DeactivateSubscription(ctx context.Context, userID int) error {
	query := `UPDATE users
			  SET subscription_active = FALSE
		      WHERE telegram_id = $1;`
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
