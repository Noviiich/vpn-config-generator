package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(username string, password string, dbName string, host string, port string, sslmode string) *Storage {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, username, dbName, password, sslmode))
	if err != nil {
		log.Fatal("Ошибка подключения к postgreSQL:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка пинга БД:", err)
	}

	return &Storage{db: db}
}
