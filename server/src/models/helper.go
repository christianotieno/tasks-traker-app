package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/christianotieno/tasks-traker-app/server/src/config"
	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func (tm *UserModel) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		log.Println("Bad Input:", err)
		return
	}

	user, err := tm.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Println("Invalid credentials:", err)
		return
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Println("Invalid credentials when checking password:", err)
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")
	// Generate JWT token
	token, err := config.GenerateToken(user.ID, string(user.Role), secret)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to generate token:", err)
		return
	}

	// Construct the response JSON
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	// Set response headers and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (tm *UserModel) generateToken(userID int, role string) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	token, tokenErr := config.GenerateToken(userID, role, secret)

	if tokenErr != nil {
		log.Println("Failed to generate token", tokenErr)
		return nil, tokenErr
	}

	return []byte(token), nil
}

func (tm *UserModel) GetUserByEmail(email string) (*entities.UserJSON, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	email = strings.ToLower(email)
	row := tm.Db.QueryRow("SELECT id, first_name, last_name, email, password, role FROM users WHERE email = ?", email)

	user := entities.UserJSON{}

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Println("Error retrieving user:", err)
		return nil, err
	}
	return &user, nil
}
