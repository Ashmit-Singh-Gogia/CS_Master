package routes

import (
	handlers "CS_Master/internal/handlers"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Quiz API Running"))
	})

	// Question routes
	r.Post("/questions", handlers.CreateQuestions(db))
	r.Get("/questions", handlers.GetAllQuestions(db))
	r.Get("/questions/{id}", handlers.GetOneQuestion(db))
	r.Put("/questions/{id}", handlers.UpdateQuestion(db))
	r.Delete("/questions/{id}", handlers.DeleteQuestion(db))
	r.Post("/questions/{id}/check", handlers.CheckAnswer(db))

	r.Post("/signup", handlers.Signup(db))
	r.Post("/login", handlers.Login(db))

	return r
}
