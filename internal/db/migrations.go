package db

import (
	"database/sql"
	"log"
)

func CreateQuestionsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS questions(
		id SERIAL PRIMARY KEY,
		questions_text TEXT NOT NULL,
		options TEXT NOT NULL,
		correct_index INT NOT NULL
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:, ", err)
	}
}

func CreateUsersTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}
}

func RunMigrations(db *sql.DB) {
	CreateQuestionsTable(db)
	CreateUsersTable(db)
}
