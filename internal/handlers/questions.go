package handlers

import (
	"CS_Master/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func CreateQuestions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q models.Question
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

func GetAllQuestions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, questions_text, options, correct_index FROM questions")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer rows.Close()
		var list []models.Question
		for rows.Next() {
			var q models.Question
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

func GetOneQuestion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var q models.Question
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

func UpdateQuestion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var q models.Question
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

func DeleteQuestion(db *sql.DB) http.HandlerFunc {
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

func CheckAnswer(db *sql.DB) http.HandlerFunc {
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
