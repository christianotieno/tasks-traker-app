package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userID from context"))
		return
	}

	userRole, ok := r.Context().Value("userRole").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userRole from context"))
		return
	}

	// Check if the user role is “Technician”
	if userRole != "Technician" {
		http.Error(w, "Only Technicians can create tasks", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		log.Println("Bad Input", err)
		return
	}
	task := entities.Task{}
	task.UserID = userID
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		log.Println("Unmarshalling failed:", err)
		return
	}

	result, err := tm.Db.Exec("INSERT INTO tasks (user_id, summary, date) VALUES (?, ?, ?)", userID, task.Summary, task.Date)
	if err != nil {
		http.Error(w, "Task creation failed", http.StatusInternalServerError)
		log.Println("Task creation failed:", err)
		return
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to retrieve created taskID:", err)
		return
	}

	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID)
	err = row.Scan(&task.ID, &task.UserID, &task.Summary, &task.Date)
	if err != nil {
		http.Error(w, "Failed to retrieve created task", http.StatusInternalServerError)
		log.Println("Failed to retrieve created task:", err)
		return
	}

	// Serialize the created task to JSON
	responseJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to serialize response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to write response:", err)
		return
	}
}

func (tm *TaskModel) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	tasksRow := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID)

	var task entities.Task
	taskErr := tasksRow.Scan(&task.ID, &task.UserID, &task.Summary, &task.Date)
	if taskErr != nil {
		if errors.Is(taskErr, sql.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Failed to scan task:", taskErr)
		}
		return
	}

	userRole, ok := r.Context().Value("userRole").(string)
	fmt.Printf("UserRole: %s\n", r.Context().Value("userRole"))
	fmt.Printf("UserID: %s\n", r.Context().Value("userID"))

	fmt.Printf("userRole: %s\n", userRole)

	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userRole from context"))
		return
	}

	// Check if the user role is “Manager”
	if userRole != "Manager" {
		http.Error(w, "Only Managers can delete tasks", http.StatusForbidden)
		return
	}

	// Delete the task from the database
	_, err = tm.Db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		http.Error(w, "Task deletion failed", http.StatusInternalServerError)
		log.Println("Task deletion failed:", err)
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

func (tm *TaskModel) UpdateTask(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	tasksRow := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID)

	var task entities.Task
	taskErr := tasksRow.Scan(&task.ID, &task.UserID, &task.Summary, &task.Date)
	if taskErr != nil {
		if errors.Is(taskErr, sql.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Failed to scan tasks:", taskErr)
		}
		return
	}

	userRole, ok := r.Context().Value("userRole").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)

		log.Println(errors.New("failed to retrieve userRole from context"))
		return
	}

	// Check if the user role is "Technician"
	if userRole != "Technician" {
		http.Error(w, "Only Technicians can update their tasks", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userID from context"))
		return
	}

	// Check if the task belongs to the user
	if task.UserID != userID {
		http.Error(w, "You can only update your own tasks", http.StatusForbidden)
		return
	}

	decoder := json.NewDecoder(r.Body)
	patchTask := make(map[string]interface{})
	err = decoder.Decode(&patchTask)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		log.Println("Decoding failed:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close request body:", err)
		}
	}(r.Body)

	// Update the task with the provided fields
	if summary, ok := patchTask["summary"].(string); ok {
		task.Summary = summary
	}

	if date, ok := patchTask["date"].(string); ok {
		task.Date = date
	}

	// Update the task in the database
	_, err = tm.Db.Exec("UPDATE tasks SET summary = ?, date = ? WHERE id = ?", task.Summary, task.Date, taskID)
	if err != nil {
		http.Error(w, "Task update failed", http.StatusInternalServerError)
		log.Println("Failed to update task:", err)
		return
	}

	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID)

	var updatedTask entities.Task
	err = row.Scan(&updatedTask.ID, &updatedTask.UserID, &updatedTask.Summary, &updatedTask.Date)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to get updated Task:", err)
		return
	}

	responseJSON, err := json.Marshal(updatedTask)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to serialize response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println("Failed to write response:", err)
		return
	}
}
