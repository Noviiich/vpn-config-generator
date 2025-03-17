package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

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

func (s *Storage) InitDB(schemaFile string) error {
	sqlBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("can't read SQL-file: %w", err)
	}

	_, err = s.pool.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		return fmt.Errorf("can't run SQL: %w", err)
	}

	return nil
}
