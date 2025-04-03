package postgres

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Storage) UpdateIP(ctx context.Context, newIP string) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO ip_pool (last_ip) 
        VALUES ($1)`, newIP)
	return err
}

func (s *Storage) GetIP(ctx context.Context) (string, error) {
	var lastIP string
	err := s.pool.QueryRow(ctx, `
        SELECT last_ip 
        FROM ip_pool 
        ORDER BY id DESC 
        LIMIT 1
    `).Scan(&lastIP)

	if err != nil {
		// Если таблица пустая — возвращаем дефолтный IP
		if errors.Is(err, sql.ErrNoRows) {
			return "10.0.0.1", nil
		}
		return "", err
	}

	return lastIP, nil
}
