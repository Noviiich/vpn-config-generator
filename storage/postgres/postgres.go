package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(username string, password string, dbName string) *Storage {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", username, password, dbName)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к postgreSQL:", err)
	}

	return &Storage{pool: pool}
}

// func (s *Storage) InitDB(schemaFile string) error {
// 	sqlBytes, err := os.ReadFile(schemaFile)
// 	if err != nil {
// 		return fmt.Errorf("can't read SQL-file: %w", err)
// 	}

// 	_, err = s.pool.Exec(context.Background(), string(sqlBytes))
// 	if err != nil {
// 		return fmt.Errorf("can't run SQL: %w", err)
// 	}

// 	return nil
// }

// func (s *Storage) initDB(schemaFile string) error {
// 	sqlBytes, err := os.ReadFile(schemaFile)
// 	if err != nil {
// 		return fmt.Errorf("can't read SQL-file: %w", err)
// 	}

// 	_, err = s.pool.Exec(context.Background(), string(sqlBytes))
// 	if err != nil {
// 		return fmt.Errorf("can't run SQL: %w", err)
// 	}

// 	return nil
// }

func (s *Storage) InitDB() error {
	initCommand := `
CREATE TABLE IF NOT EXISTS users (
telegram_id BIGINT PRIMARY KEY,
username TEXT UNIQUE NOT NULL,
subscription_active BOOLEAN DEFAULT FALSE,
subscription_expiry TIMESTAMP
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id) ON DELETE CASCADE,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    ip TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);`

	_, err := s.pool.Exec(context.Background(), initCommand)
	if err != nil {
		return fmt.Errorf("can't init SQL: %w", err)
	}

	return nil
}
