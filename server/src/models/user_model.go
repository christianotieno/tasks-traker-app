package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
)

type UserModel struct {
	Db *sql.DB
}

func UserHandler(db *sql.DB) *UserModel {
	return &UserModel{
		Db: db,
	}
}

func (tm *UserModel) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	user := entities.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	// Validate the role
	if user.Role != "Manager" && user.Role != "Technician" {
		http.Error(w, "Invalid role, try again!", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	// Hash the user’s password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Insert the user details into the database
	result, err := tm.Db.Exec(
		"INSERT INTO users (first_name, last_name, email, role, password) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Role, string(hashedPassword))
	if err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Retrieve the ID of the created user
	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &entities.JWTClaims{
		UserID: int(userID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	// Sign the token with the JWT signing key
	jwtKey := []byte(secret)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Clear the password field in the user struct before serializing to JSON
	user.Password = ""

	// Serialize the created user account details and the JWT token to JSON
	responseJSON, err := json.Marshal(struct {
		User  entities.User `json:"user"`
		Token string        `json:"token"`
	}{
		User:  user,
		Token: tokenString,
	})
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}

func (tm *UserModel) GetAllTasksByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Retrieve the user ID from the request context
	userID := r.Context().Value("userID")
	fmt.Printf("userID: %v\n", userID)

	// Only allow access to tasks belonging to the user
	rows, err := tm.Db.Query("SELECT id, summary, date FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving data", err)
		return
	}

	fmt.Printf("rows: %v\n", rows)

	// Extract the JWT token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	// Verify and parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		return
	}

	// Check if the token is valid and has not expired
	if !token.Valid {
		http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}(rows)

	var tasks []entities.Task

	for rows.Next() {
		task := entities.Task{}
		err := rows.Scan(&task.ID, &task.Summary, &task.Date)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		tasks = append(tasks, task)
	}

	response, err := json.Marshal(&tasks)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}

func (tm *UserModel) GetAllUsersAndAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	userID := r.Context().Value("userID").(int)
	isManager := tm.isManager(userID)
	if !isManager {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Unauthorized access")
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	secret := os.Getenv("SECRET")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		return
	}

	userRows, usersErr := tm.Db.Query("SELECT id, firstname, lastname, email, role FROM users")
	if usersErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(usersErr)
		return
	}

	defer func(rows *sql.Rows) {
		err := userRows.Close()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}(userRows)

	var users []entities.User

	for userRows.Next() {
		user := entities.User{}
		err := userRows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		users = append(users, user)
	}

	for i := range users {
		rows, userErr := tm.Db.Query("SELECT id, summary, date FROM tasks WHERE user_id = ?", users[i].ID)
		if userErr != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(userErr)
			return
		}

		var tasks []entities.Task

		for rows.Next() {
			task := entities.Task{}
			err := rows.Scan(&task.ID, &task.Summary, &task.Date)
			if err != nil {
				err := rows.Close()
				if err != nil {
					return
				}
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			tasks = append(tasks, task)
		}

		users[i].Tasks = &tasks

		if err := rows.Close(); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	response, resErr := json.Marshal(&users)
	if resErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(resErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(response)
	if writeErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(writeErr)
		return
	}
}

func (tm *UserModel) GetUserByEmail(email string) (*entities.User, error) {
	row := tm.Db.QueryRow("SELECT id, first_name, last_name, email, password, role FROM users WHERE email = ?", email)

	user := entities.User{}

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

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
