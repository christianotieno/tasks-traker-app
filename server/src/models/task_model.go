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

type TaskModel struct {
	Db *sql.DB
}

func TaskHandler(db *sql.DB) *TaskModel {
	return &TaskModel{
		Db: db,
	}
}

func (tm *TaskModel) CreateTask(w http.ResponseWriter, r *http.Request, technicianID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		return
	}

	var task entities.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	// Insert the task into the database
	result, err := tm.Db.Exec("INSERT INTO tasks (summary, date, technician_id) VALUES (?, ?)", task.Summary, task.Date, technicianID)
	if err != nil {
		http.Error(w, "Task creation failed", http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the created task
	taskID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve task ID", http.StatusInternalServerError)
		return
	}

	// Retrieve the created task from the database
	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID)
	err = row.Scan(&task.ID, &task.Summary, &task.Date)
	if err != nil {
		http.Error(w, "Failed to retrieve created task", http.StatusInternalServerError)
		return
	}

	// Serialize the created task to JSON
	responseJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (tm *TaskModel) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := tm.Db.Query("SELECT * FROM tasks")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
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

	// Retrieve tasks from the database
	for rows.Next() {
		task := entities.Task{}
		err := rows.Scan(&task.ID, &task.Summary, &task.Date)
		if err != nil {
			http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		tasks = append(tasks, task)
	}

	// Serialize tasks to JSON
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

func (tm *TaskModel) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task entities.Task
	err := row.Scan(&task.ID, &task.Summary, &task.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
		}
		return
	}

	// Serialize task to JSON
	response, err := json.Marshal(&task)
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

func (tm *TaskModel) UpdateTask(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body into a Task struct
	var updatedTask entities.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the task exists
	var count int
	err = tm.Db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ?", id).Scan(&count)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	if count == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Update the task in the database
	_, err = tm.Db.Exec("UPDATE tasks SET summary = ?, date = ? WHERE id = ?", updatedTask.Summary, updatedTask.Date, id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Serialize the updated task to JSON
	responseJSON, err := json.Marshal(updatedTask)
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
			log.Fatal(taskErr)
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
			log.Fatal(userErr)
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
