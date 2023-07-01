package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Task struct {
	Summary string    `json:"summary"`
	Date    time.Time `json:"date"`
}

type TaskHandler struct {
	db *sql.DB
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{
		db: db,
	}
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Insert the task into the database
	_, err = th.db.Exec("INSERT INTO tasks (summary, date) VALUES (?, ?)", task.Summary, task.Date)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
