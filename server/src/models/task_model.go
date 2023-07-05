package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
)

type TaskModel struct {
	Db *sql.DB
}

func TaskHandler(db *sql.DB) *TaskModel {
	return &TaskModel{
		Db: db,
	}
}

func (tm *TaskModel) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	userID := r.Context().Value("userID")
	fmt.Printf("userID: %+v\n", userID)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		log.Println("Bad Input", err)
		return
	}
	task := entities.Task{}
	task.UserID = userID.(int)
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		log.Println("Unmarshalling failed", err)
		return
	}

	// Insert the task into the database
	result, err := tm.Db.Exec("INSERT INTO tasks (summary, date, user_id) VALUES (?, ?, ?)", task.Summary, task.Date, userID)
	if err != nil {
		fmt.Printf("Task: %+v\n", task)
		http.Error(w, "Task creation failed", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Retrieve the ID of the created task
	taskID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve task ID", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Retrieve the created task from the database
	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID)
	err = row.Scan(&task.ID, &task.Summary, &task.Date, &task.UserID)
	if err != nil {
		http.Error(w, "Failed to retrieve created task", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Serialize the created task to JSON
	responseJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (tm *TaskModel) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasksRow := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task entities.Task
	taskErr := tasksRow.Scan(&task.ID, &task.Summary, &task.Date, &task.UserID)
	if taskErr != nil {
		if errors.Is(taskErr, sql.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(taskErr)
		}
		return
	}

	usersRow := tm.Db.QueryRow("SELECT * FROM users WHERE id = ?", task.UserID)

	var user entities.User
	userErr := usersRow.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
	if userErr != nil {
		if errors.Is(userErr, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println(userErr)
		}
		return
	}

	// Check if the user role is “Manager”
	if user.Role != "Manager" {
		http.Error(w, "Only Managers can delete tasks", http.StatusForbidden)
		return
	}

	// Delete the task from the database
	_, err := tm.Db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Task deletion failed", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Task deleted successfully",
	}

	// Serialize the response to JSON
	responseJSON, err := json.Marshal(response)
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
