package database

import (
	"github.com/jmoiron/sqlx"
)

func CreateTableUsers(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL
	);
	`)
	return err
}

func CreateTableMessages(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS messages (
		id INT PRIMARY KEY,
		sender_username TEXT NOT NULL,
		receiver_username TEXT NOT NULL,
		body TEXT NOT NULL,
		sent_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`)
	return err
}
