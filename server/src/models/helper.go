package models

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	userID := user.ID
	managerID := user.ManagerID
	token, err := config.GenerateToken(userID, managerID, secret)
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

func (tm *UserModel) generateToken(userID string, managerID string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	token, tokenErr := config.GenerateToken(userID, managerID, secret)

	if tokenErr != nil {
		log.Println("Failed to generate token", tokenErr)
		return "", tokenErr
	}

	return token, nil
}

func (tm *UserModel) GetUserByEmail(email string) (*entities.UserJSON, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	email = strings.ToLower(email)

	user := entities.UserJSON{}

	usersRow := tm.Db.QueryRow("SELECT id, first_name, last_name, email, password FROM users WHERE email = ?", email)
	err := usersRow.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Println("Error retrieving user:", err)
		return nil, err
	}

	managersRow := tm.Db.QueryRow("SELECT manager_id FROM managers WHERE technician_id = ?", user.ID)
	var managerID string
	err = managersRow.Scan(&managerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		log.Println("Error retrieving manager:", err)
		return nil, err
	}

	user.ManagerID = managerID

	fmt.Printf("User: %+v\n", user)
	return &user, nil
}

func (tm *UserModel) validateUser(user entities.UserJSON, w http.ResponseWriter) error {
	if user.FirstName == "" {
		return httpError(w, http.StatusBadRequest, "Missing required fields: first_name")
	}

	if user.LastName == "" {
		return httpError(w, http.StatusBadRequest, "Missing required fields: last_name")
	}

	if user.Email == "" {
		return httpError(w, http.StatusBadRequest, "Missing required fields: email")
	}

	if user.Password == "" {
		return httpError(w, http.StatusBadRequest, "Missing password")
	}

	if len(user.Password) < 6 {
		return httpError(w, http.StatusBadRequest, "Password must be at least 6 characters")
	}

	if !strings.Contains(user.Email, "@") {
		return httpError(w, http.StatusBadRequest, "Invalid email address")
	}

	if tm.userExists(user.Email) {
		return httpError(w, http.StatusBadRequest, "Email already exists, please try again with a different email")
	}

	return nil
}

func httpError(w http.ResponseWriter, statusCode int, message string) error {
	http.Error(w, message, statusCode)
	log.Println(message)
	return errors.New(message)
}
