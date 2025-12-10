package handlers

import (
	"CS_Master/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var user models.User
		err = db.QueryRow(`
			SELECT id, username, password_hash FROM users WHERE username=$1
		`, creds.UserName).Scan(&user.ID, &user.Username, &user.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Step 4: Create JWT Token
		secret := []byte("SUPER_SECRET_KEY") // will move to env later

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(secret)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}
		// Step 5: Send token to user
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})

	}
}
