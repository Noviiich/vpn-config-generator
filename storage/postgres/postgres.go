package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(username string, password string, dbName string) *Storage {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"localhost", "5432", username, dbName, password, "disable"))
	if err != nil {
		log.Fatal("Ошибка подключения к postgreSQL:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка пинга БД:", err)
	}

	return &Storage{db: db}
}

func (s *Storage) InitDB(ctx context.Context) {
	query := `
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

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("can't init postgeSQL: %v", err)
	}
}
