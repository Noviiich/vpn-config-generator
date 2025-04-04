package postgres

import (
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
