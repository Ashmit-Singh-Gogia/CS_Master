package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

type Question struct {
	ID           int      `json:"id"`
	QuestionText string   `json:"questions_text"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correct_index"`
}

func main() {
	connStr := "host=localhost port=5433 user=postgres password=postgres dbname=CS_MasterDB sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Can't connect to DB: ", err)
	}
	fmt.Println("Connected to CS_MasterDB successfully!")
	createQuestionsTable(db)
	fmt.Println("Questions table ready!")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Quiz Api Running"))
	})
	r.Post("/questions", CreateQuestions(db))
	r.Get("/questions", getAllQuestions(db))
	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}

func createQuestionsTable(db *sql.DB) {
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

func CreateQuestions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q Question
		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		// convert the slice of string options into json string
		opts, err := json.Marshal(q.Options)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		err = db.QueryRow(
			"INSERT INTO questions (questions_text, options, correct_index) VALUES ($1, $2, $3) RETURNING id",
			q.QuestionText, string(opts), q.CorrectIndex,
		).Scan(&q.ID)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode("Question Created Successfully")
	}
}

func getAllQuestions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, questions_text, options, correct_index FROM questions")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer rows.Close()
		var list []Question
		for rows.Next() {
			var q Question
			var opts string

			err := rows.Scan(&q.ID, &q.QuestionText, &opts, &q.CorrectIndex)
			if err != nil {
				continue
			}
			json.Unmarshal([]byte(opts), &q.Options)

			list = append(list, q)
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", " ")
		enc.Encode(list)
	}
}
