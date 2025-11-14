package database

import (
	"fmt"
	"messenger-project/pkg/config"
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitPostgres() (*sqlx.DB, error) {
	cfg := config.Load()

	totalTimeout := 30 * time.Second
	backoff := 500 * time.Millisecond
	deadline := time.Now().Add(totalTimeout)

	var db *sqlx.DB
	var err error
	for {
		db, err = sqlx.Connect("postgres", cfg.DATABASE_URL)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				DB = db
				return db, nil
			} else {
				_ = db.Close()
				err = pingErr
			}
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("failed to connect to Postgres after %s: %w", totalTimeout.String(), err)
		}

		time.Sleep(backoff)
		if backoff < 5*time.Second {
			backoff = backoff * 2
		}
	}
}
