package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
		log.Println(err)
		return
	}

	user := entities.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		log.Println("Invalid JSON format:", err)
		return
	}

	// Check if the email already exists
	if tm.userExists(user.Email) {
		http.Error(w, "Email already exists, please try again with a different email", http.StatusBadRequest)
		return
	}

	// Validate the role
	if !tm.isValidRole(string(user.Role)) {
		http.Error(w, "Invalid role, try again!", http.StatusBadRequest)
		return
	}

	// Hash user password
	hashedPassword, err := tm.hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Failed to hash password:", err)
		return
	}

	// Insert user details into the database
	userID, err := tm.insertUser(user.FirstName, user.LastName, user.Email, string(hashedPassword), string(user.Role))
	if err != nil {
		http.Error(w, "Something went wrong, try again", http.StatusInternalServerError)
		log.Fatal("User creation failed:", err)
		return
	}

	// Generate JWT token
	tokenString, err := tm.generateToken(int(userID))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Failed to generate token:", err)
		return
	}

	user.Password = ""

	responseJSON, err := json.Marshal(struct {
		User  entities.User `json:"user"`
		Token string        `json:"token"`
	}{
		User: entities.User{
			ID:        int(userID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			Tasks:     user.Tasks,
		},
		Token: string(tokenString),
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Failed to serialize responseJSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Failed to write response:", err)
		return
	}
}

func (tm *UserModel) userExists(email string) bool {
	var count int
	err := tm.Db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Println("Failed to query user database:", err)
		return true
	}
	return count > 0
}

func (tm *UserModel) isValidRole(role string) bool {
	return role == "Manager" || role == "Technician"
}

func (tm *UserModel) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (tm *UserModel) insertUser(firstName, lastName, email, password, role string) (int64, error) {
	result, err := tm.Db.Exec("INSERT INTO users (first_name, last_name, email, password, role) VALUES (?, ?, ?, ?, ?)", firstName, lastName, email, password, role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (tm *UserModel) GetAllTasksByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Retrieve the user ID from the request context
	userID := r.Context().Value("userID")

	// Only allow access to tasks belonging to the user
	rows, err := tm.Db.Query("SELECT id, summary, date FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving data", err)
		return
	}

	// Extract the JWT token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

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

	response, err := json.Marshal(tasks)
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

	userIDFloat := r.Context().Value("userID").(float64)
	userID := int(userIDFloat)

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

	userRows, usersErr := tm.Db.Query("SELECT id, first_name, last_name, email, role FROM users")
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
