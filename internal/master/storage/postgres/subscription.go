package postgres

import (
	"context"
	"database/sql"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

func (s *Storage) CreateSubscription(ctx context.Context, userID, typeID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO subscriptions (user_id, type_id, expiry_date)
		VALUES ($1, $2, 
		NOW() + (SELECT duration FROM subscription_types WHERE id = $2)
		)
		ON CONFLICT (user_id) 
		DO UPDATE SET
			type_id = EXCLUDED.type_id,
			expiry_date = NOW() + (
				SELECT duration FROM subscription_types 
				WHERE id = EXCLUDED.type_id
			)`

	_, err = tx.ExecContext(ctx, query, userID, typeID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DeleteSubscription(ctx context.Context, userID int) error {
	query := `
			DELETE FROM subscriptions 
			WHERE user_id = $1`
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *Storage) GetSubscription(ctx context.Context, userID int) (*storage.Subscription, error) {
	query := `
			SELECT *
			FROM subscriptions
			WHERE user_id = $1
			LIMIT 1`

	var sub storage.Subscription
	err := s.db.GetContext(ctx, &sub, query, userID)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *Storage) GetSubscriptions(ctx context.Context, userID int) ([]storage.Subscription, error) {
	query := `
			SELECT *
			FROM subscriptions`

	var subs []storage.Subscription
	err := s.db.SelectContext(ctx, &subs, query)
	if err != nil {
		return nil, err
	}

	if subs == nil {
		return nil, sql.ErrNoRows
	}

	return subs, nil
}
