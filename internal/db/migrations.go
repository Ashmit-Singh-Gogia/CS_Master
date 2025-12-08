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
