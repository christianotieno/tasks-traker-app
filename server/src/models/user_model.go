package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

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

	var user entities.User
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

	// Insert the user details into the database
	result, err := tm.Db.Exec(
		"INSERT INTO users (first_name, last_name, email, role) VALUES (?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Role)
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

	// Retrieve the created user account from the database
	row := tm.Db.QueryRow("SELECT * FROM users WHERE id = ?", userID)
	err = row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
	if err != nil {
		http.Error(w, "Failed to retrieve created user", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Serialize the created user account details to JSON
	responseJSON, err := json.Marshal(user)
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

func (tm *UserModel) GetUser(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	row := tm.Db.QueryRow("SELECT id, firstname, lastname, email, role FROM users WHERE id = ?", userID)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		return
	}

	response, err := json.Marshal(user)
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

func (tm *UserModel) GetAllTasksByUserID(w http.ResponseWriter, r *http.Request, technicianID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	rows, err := tm.Db.Query("SELECT id, summary, date FROM tasks WHERE user_id = ?", technicianID)
	if err != nil {
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}(rows)

	var tasks []entities.Task

	for rows.Next() {
		task := entities.Task{}
		err := rows.Scan(&task.ID, &task.Summary, &task.Date)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
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

func (tm *UserModel) GetAllUsersAndAllTasks(w http.ResponseWriter, r *http.Request, managerID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	userRows, usersErr := tm.Db.Query("SELECT id, firstname, lastname, email, role FROM users")
	if usersErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(usersErr)
		return
	}

	defer func(rows *sql.Rows) {
		err := userRows.Close()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}(userRows)

	var users []entities.User

	for userRows.Next() {
		user := entities.User{}
		err := userRows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		users = append(users, user)
	}

	for i := range users {
		rows, userErr := tm.Db.Query("SELECT id, summary, date FROM tasks WHERE user_id = ?", users[i].ID)
		if userErr != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(userErr)
			return
		}

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Fatal(err)
				return
			}
		}(rows)

		var tasks []entities.Task

		for rows.Next() {
			task := entities.Task{}
			err := rows.Scan(&task.ID, &task.Summary, &task.Date)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				log.Fatal(err)
				return
			}
			tasks = append(tasks, task)
		}

		users[i].Tasks = tasks
	}

	response, resErr := json.Marshal(&users)
	if resErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(resErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(response)
	if writeErr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(writeErr)
		return
	}
}
