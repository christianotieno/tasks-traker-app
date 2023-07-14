package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
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
		log.Println("Bad Input", err)
		return
	}

	user := entities.UserJSON{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		log.Println("Invalid JSON format:", err)
		return
	}

	err = tm.validateUser(user, w)
	if err != nil {
		return
	}

	// Hash user password
	hashedPassword, err := tm.hashPassword([]byte(user.Password))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Failed to hash password:", err)
		return
	}

	// Insert user details into the database
	userID, err := tm.insertUser(user.FirstName, user.LastName, user.Email, hashedPassword, user.ManagerID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Account creation failed:", err)
		return
	}

	// Generate JWT token
	tokenString, err := tm.generateToken(userID, user.ManagerID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal("Failed to generate token:", err)
		return
	}

	responseJSON, err := json.Marshal(struct {
		User    entities.User `json:"user"`
		Message string        `json:"message"`
		Token   string        `json:"token"`
	}{
		User: entities.User{
			ID:        userID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			ManagerID: user.ManagerID,
			Tasks:     user.Tasks,
		},
		Message: "Account creation successful",
		Token:   tokenString,
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

func (tm *UserModel) hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (tm *UserModel) insertUser(firstName, lastName, email string, hashedPassword []byte, managerID string) (string, error) {
	userID := uuid.New().String()

	// Insert the technician or manager
	query := "INSERT INTO users (id, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)"
	_, err := tm.Db.Exec(query, userID, firstName, lastName, email, hashedPassword)
	if err != nil {
		return "", err
	}

	if managerID != "" {
		// Check if the manager exists
		if err != nil {
			return "", fmt.Errorf("invalid managerID")
		}
		managerExistsQuery := "SELECT id FROM users WHERE id = ?"
		row := tm.Db.QueryRow(managerExistsQuery, managerID)
		var mID string
		err = row.Scan(&mID)
		fmt.Printf("managerID: %v\n", managerID)
		if err != nil {
			return "", fmt.Errorf("manager does not exist")
		}

		// Insert the manager-technician relationship if managerID is not nil
		if managerID != "" {
			id := uuid.New().String()
			relationshipQuery := "INSERT INTO managers (id, manager_id, technician_id) VALUES (?, ?, ?)"
			_, err = tm.Db.Exec(relationshipQuery, id, managerID, userID)
			if err != nil {
				return "", err
			}
		}
	}

	return userID, nil
}

func (tm *UserModel) GetAllTasksByUserID(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Retrieve the user ID from the request context
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to retrieve user ID")
		return
	}

	managerID, ok := r.Context().Value("managerID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to retrieve user role")
		return
	}

	// Extract the JWT token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	// Fetch tasks from the database based on the user ID
	rows, err := tm.Db.Query("SELECT id, summary, date, user_id FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Error retrieving data", err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}(rows)

	var tasks []entities.Task

	for rows.Next() {
		task := entities.Task{}
		err := rows.Scan(&task.ID, &task.Summary, &task.Date, &task.UserID)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Only allow access if the user is a manager or the technician of the specified user ID
		if !(id == managerID || id == task.UserID) {
			http.Error(w, "Access denied", http.StatusForbidden)
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

	// Retrieve the user role from the request context
	userRole, ok := r.Context().Value("userRole").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to retrieve user role")
		return
	}

	// Only allow access if the user is a manager
	if userRole != "Manager" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Unauthorized access")
		return
	}

	userRows, usersErr := tm.Db.Query("SELECT id, first_name, last_name, email, role, manager_id FROM users")
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
		err := userRows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.ManagerID)
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
