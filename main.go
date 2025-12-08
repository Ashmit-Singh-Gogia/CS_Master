package main

import (
	dbpkg "CS_Master/internal/db"
	handlers "CS_Master/internal/handlers"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	_ "github.com/lib/pq"
)

func main() {

	db := dbpkg.Connect()
	fmt.Println("Connected to CS_MasterDB successfully!")
	dbpkg.CreateQuestionsTable(db)
	fmt.Println("Questions table ready!")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Quiz Api Running"))
	})
	r.Post("/questions", handlers.CreateQuestions(db))
	r.Get("/questions", handlers.GetAllQuestions(db))
	r.Get("/questions/{id}", handlers.GetOneQuestion(db))
	r.Put("/questions/{id}", handlers.UpdateQuestion(db))
	r.Delete("/questions/{id}", handlers.DeleteQuestion(db))
	r.Post("/questions/{id}/check", handlers.CheckAnswer(db))
	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
