package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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
		return
	}

	user, err := tm.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Printf("user: %v\n", user)

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	token, tokenErr := GenerateJWTToken(user.ID, secret)

	if tokenErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to generate token", tokenErr)
		return
	}

	responseJSON, err := json.Marshal(struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (tm *UserModel) isManager(userID int) bool {
	var role string
	err := tm.Db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		log.Println("Error retrieving user role:", err)
		return false
	}
	return role == "Manager"
}

func (tm *UserModel) generateToken(i int) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	token, tokenErr := GenerateJWTToken(i, secret)

	if tokenErr != nil {
		log.Println("Failed to generate token", tokenErr)
		return nil, tokenErr
	}

	return []byte(token), nil
}

func GenerateJWTToken(userID int, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (1 day)
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
