package postgres

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) UpdateIP(ctx context.Context, newIP string) error {
	query := `INSERT INTO ip_pool (last_ip) VALUES ($1)`
	_, err := s.db.ExecContext(ctx, query, newIP)
	return err
}

func (s *Storage) GetIP(ctx context.Context) (string, error) {
	query := `
        SELECT last_ip 
        FROM ip_pool 
        ORDER BY id DESC 
        LIMIT 1`
	var lastIP string
	err := s.db.GetContext(ctx, &lastIP, query)

	if err != nil {
		// Если таблица пустая — возвращаем дефолтный IP
		if errors.Is(err, sql.ErrNoRows) {
			return "10.0.0.1", nil
		}
		return "", err
	}

	return lastIP, nil
}
