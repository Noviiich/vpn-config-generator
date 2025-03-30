package postgres

import (
	"context"
	"fmt"

	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *Storage) ActivateSubscription(ctx context.Context, userID, typeID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO subscriptions (user_id, type_id, is_active, expiry_date)
		VALUES ($1, $2, TRUE,
		NOW() + (SELECT duration FROM subscription_types WHERE id = $2)
		)
		ON CONFLICT (user_id) 
		DO UPDATE SET
			type_id = EXCLUDED.type_id,
			start_date = NOW(),
			expiry_date = NOW() + (
				SELECT duration FROM subscription_types 
				WHERE id = EXCLUDED.type_id
			),
			is_active = TRUE`

	_, err = tx.ExecContext(ctx, query, userID, typeID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DeactivateSubscription(ctx context.Context, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `UPDATE subscriptions
			  SET 
			  	is_active = FALSE,
				expiry_date = GREATEST(NOW(), expiry_date)
		      WHERE user_id = $1;`

	result, err := tx.ExecContext(ctx, query, userID)
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

	return tx.Commit()
}

func (s *Storage) CheckSubscription(ctx context.Context, userID int) (*storage.Subscription, error) {
	query := `
			SELECT *
			FROM subscriptions
			WHERE user_id = $1`

	var sub storage.Subscription
	err := s.db.GetContext(ctx, &sub, query, userID)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}
