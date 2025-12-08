package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	r.Get("/questions/{id}", getOneQuestion(db))
	r.Put("/questions/{id}", updateQuestion(db))
	r.Delete("/questions/{id}", deleteQuestion(db))
	r.Post("/questions/{id}/check", checkAnswer(db))
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
func getOneQuestion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var q Question
		var opts string

		err = db.QueryRow(
			"SELECT id, questions_text, options, correct_index FROM questions WHERE id=$1",
			id,
		).Scan(&q.ID, &q.QuestionText, &opts, &q.CorrectIndex)

		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.Unmarshal([]byte(opts), &q.Options)
		json.NewEncoder(w).Encode(q)
	}
}

func updateQuestion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var q Question
		err = json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		opts, err := json.Marshal(q.Options)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		result, err := db.Exec(
			"UPDATE questions SET questions_text=$1, options=$2, correct_index=$3 WHERE id=$4",
			q.QuestionText, string(opts), q.CorrectIndex, id,
		)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		numberOfRows, _ := result.RowsAffected()
		if numberOfRows == 0 {
			http.Error(w, "Question not found", 404)
			return
		}
		w.WriteHeader(204)
	}
}

func deleteQuestion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		result, err := db.Exec("DELETE FROM questions WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		rows, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if rows == 0 {
			http.Error(w, "Not Found or Already Deleted", 404)
			return
		}
		w.WriteHeader(204)
	}
}

func checkAnswer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var AnswerBody struct {
			Answer int `json:"answer"`
		}
		err = json.NewDecoder(r.Body).Decode(&AnswerBody)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var correctIndex int

		err = db.QueryRow(
			"SELECT correct_index FROM questions WHERE id = $1",
			id,
		).Scan(&correctIndex)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		isCorrect := AnswerBody.Answer == correctIndex
		json.NewEncoder(w).Encode(map[string]bool{
			"correct": isCorrect,
		})
	}
}
