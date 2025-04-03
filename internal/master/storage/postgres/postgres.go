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

func (s *Storage) InitDB(ctx context.Context) {
	initCommand := `
CREATE TABLE IF NOT EXISTS users (
telegram_id BIGINT PRIMARY KEY,
username TEXT UNIQUE NOT NULL,
subscription_active BOOLEAN DEFAULT FALSE,
subscription_expiry TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id) ON DELETE CASCADE,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    ip TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS ip_pool (
	id SERIAL PRIMARY KEY,
	last_ip TEXT NOT NULL
);`

	_, err := s.pool.Exec(ctx, initCommand)
	if err != nil {
		log.Fatalf("can't init postgeSQL: %v", err)
	}
}
