package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/storage"
)

// func (s *Storage) GetIP(ctx context.Context, deviceID int) error {
// 	query := `SELECT INTO ip_pool (last_ip) VALUES ($1)`
// 	_, err := s.db.ExecContext(ctx, query, newIP)
// 	return err
// }

// func (s *Storage) UpdateIP(ctx context.Context, newIP string) error {
// 	query := `INSERT INTO ip_pool (last_ip) VALUES ($1)`
// 	_, err := s.db.ExecContext(ctx, query, newIP)
// 	return err
// }

// func (s *Storage) CreateIP(ctx context.Context) (string, error) {
// 	query := `SELECT last_ip FROM ip_pool ORDER BY id DESC LIMIT 1`

// 	var lastIP string
// 	err := s.db.QueryRowContext(ctx, query).Scan(&lastIP)

// 	if err != nil {
// 		// Если таблица пустая — возвращаем дефолтный IP
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return "10.0.0.1", nil
// 		}
// 		return "", err
// 	}

// 	return lastIP, nil
// }

func (s *Storage) GetIpIsNull(ctx context.Context) (*storage.IPAddress, error) {
	query := `
			SELECT * FROM ip_pool
			WHERE device_id IS NULL
			LIMIT 1`
	var ip storage.IPAddress
	err := s.db.GetContext(ctx, &ip, query)
	if err != nil {
		return nil, err
	}

	return &ip, nil
}

func (s *Storage) CreateIP(ctx context.Context, IPAdderss *storage.IPAddress) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// Откат транзакции в случае ошибки
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `
			INSERT INTO ip_pool (device_id, ip, is_available)
			VALUES ($1, $2, true)
			RETURNING id`
	err = tx.QueryRowContext(ctx, query, IPAdderss.DeviceID, IPAdderss.IP).Scan(&IPAdderss.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) UpdateIpPool(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// Откат транзакции в случае ошибки
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
			UPDATE ip_pool 
			SET device_id = $1
			WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	return tx.Commit()
}
