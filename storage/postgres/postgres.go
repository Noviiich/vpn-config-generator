package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Noviiich/vpn-config-generator/storage"
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

func (s *Storage) InitDB() {
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
);

CREATE TABLE IF NOT EXISTS ip_pool (
	id SERIAL PRIMARY KEY,
	last_ip TEXT NOT NULL
);`

	_, err := s.pool.Exec(context.Background(), initCommand)
	if err != nil {
		log.Fatalf("can't init postgeSQL: %v", err)
	}
}

func (s *Storage) GetDevice(ctx context.Context, telegramID int) (storage.Device, error) {
	query := `SELECT id, user_id, private_key, public_key, ip, is_active 
              FROM devices 
              WHERE user_id = $1`

	var device storage.Device
	err := s.pool.QueryRow(ctx, query, telegramID).Scan(
		&device.ID,
		&device.UserID,
		&device.PrivateKey,
		&device.PublicKey,
		&device.IP,
		&device.IsActive,
	)

	if err != nil {
		return storage.Device{}, err
	}

	return device, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	_, err := s.pool.Exec(context.Background(),
		`INSERT INTO users (telegram_id, username, subscription_active, subscription_expiry) 
		 VALUES ($1, $2, $3, $4) ON CONFLICT (telegram_id) DO NOTHING`,
		user.TelegramID, user.Username, user.SubscriptionActive, user.SubscriptionExpiry)
	return err
}

// func (s *Storage) GetUser(ctx context.Context, telegramID int) (*storage.User, error) {
// 	query := `SELECT telegram_id, user_id, private_key, public_key, ip, is_active
//               FROM devices
//               WHERE user_id = $1`
// 	_, err := s.pool.Exec(context.Background(),
// 		`INSERT INTO users (telegram_id, username, subscription_active, subscription_expiry)
// 		 VALUES ($1, $2, $3, $4) ON CONFLICT (telegram_id) DO NOTHING`,
// 		user.TelegramID, user.Username, user.SubscriptionActive, user.SubscriptionExpiry)
// 	return err
// }

func (s *Storage) CreateDevice(ctx context.Context, device *storage.Device) error {
	query := `INSERT INTO devices (user_id, private_key, public_key, ip, is_active) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := s.pool.QueryRow(ctx, query,
		device.UserID, device.PrivateKey, device.PublicKey, device.IP, device.IsActive).Scan(&device.ID)

	return err
}

func (s *Storage) IsExistsUser(ctx context.Context, telegramID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)`
	err := s.pool.QueryRow(ctx, query, telegramID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("can't check if user exists: %v", err)
	}
	return exists, nil
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

func (s *Storage) UpdateIP(ctx context.Context, newIP string) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO ip_pool (last_ip) 
        VALUES ($1)`, newIP)
	return err
}
