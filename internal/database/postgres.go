package database

import (
	"messenger-project/pkg/config"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitPostgres() (*sqlx.DB, error) {
	cfg := config.Load()
	db, err := sqlx.Connect("postgres", cfg.DATABASE_URL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}
