package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/services"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
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

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userID from context"))
		return
	}

	managerID, ok := r.Context().Value("managerID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve managerID from context"))
		return
	}

	// Check if the user is a "Technician"
	if !(managerID != "") {
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

	id := uuid.New().String()
	// Save the task to the database
	insertQuery := "INSERT INTO tasks (id, summary, date, user_id) VALUES (?, ?, ?, ?)"
	_, err = tm.Db.Exec(insertQuery, id, task.Summary, task.Date, userID)
	if err != nil {
		http.Error(w, "Task creation failed", http.StatusInternalServerError)
		log.Println("Task creation failed:", err)
		return
	}

	kafkaProducer, err := services.NewKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Println("Failed to initialize Kafka producer:", err)
	} else {
		go func() {
			message := "New task created: " + task.Summary
			err = kafkaProducer.SendMessage([]byte(message))
			if err != nil {
				log.Println("Failed to send Kafka message:", err)
			}
		}()
	}

	// Retrieve the created task
	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	err = row.Scan(&task.ID, &task.Summary, &task.Date, &task.UserID)
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

	tasksRow := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task entities.Task
	taskErr := tasksRow.Scan(&task.ID, &task.Summary, &task.Date, &task.UserID)
	if taskErr != nil {
		if errors.Is(taskErr, sql.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Failed to scan task:", taskErr)
		}
		return
	}

	managerID, ok := r.Context().Value("managerID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve managerID from context"))
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userID from context"))
		return
	}

	// fetch managerID from the database
	managerRow := tm.Db.QueryRow("SELECT manager_id FROM managers WHERE technician_id = ?", task.UserID)
	var mID string
	managerErr := managerRow.Scan(&mID)
	if managerErr != nil {
		if errors.Is(managerErr, sql.ErrNoRows) {
			http.Error(w, "Manager not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}

	// Check if the user is a "Manager" and is manager of the user who created the task
	if !(managerID == "" && userID == mID) {
		http.Error(w, "Only Managers associated with this user can delete this task", http.StatusForbidden)
		return
	}

	// Delete the task from the database
	_, err := tm.Db.Exec("DELETE FROM tasks WHERE id = ?", id)
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

	tasksRow := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var task entities.Task
	taskErr := tasksRow.Scan(&task.ID, &task.Summary, &task.Date, &task.UserID)
	if taskErr != nil {
		if errors.Is(taskErr, sql.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Println("Failed to scan tasks:", taskErr)
		}
		return
	}

	managerID, ok := r.Context().Value("managerID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve managerID from context"))
		return
	}

	// Check if the user is a "Technician"
	if !(managerID != "") {
		http.Error(w, "Only Technicians can update their tasks", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Println(errors.New("failed to retrieve userID from context"))
		return
	}

	// Check if the task belongs to the user
	if !(task.UserID == userID) {
		http.Error(w, "Only the task owner can update this task", http.StatusForbidden)
		return
	}

	decoder := json.NewDecoder(r.Body)
	patchTask := make(map[string]interface{})
	err := decoder.Decode(&patchTask)
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
	_, err = tm.Db.Exec("UPDATE tasks SET summary = ?, date = ? WHERE id = ?", task.Summary, task.Date, id)
	if err != nil {
		http.Error(w, "Task update failed", http.StatusInternalServerError)
		log.Println("Failed to update task:", err)
		return
	}

	row := tm.Db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

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
