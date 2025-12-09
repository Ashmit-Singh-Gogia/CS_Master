package handlers

import (
	"CS_Master/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate input // i.e. check if input is valid or not
		if user.Username == "" || user.Password == "" {
			http.Error(w, "Invalid Username or Password", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.QueryRow(`
			INSERT INTO users (username, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`, user.Username, string(hash)).Scan(&user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		user.Password = "#Some PassWord#"
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Signup-successfull",
			"user":    user.Username,
		})
	}
}
