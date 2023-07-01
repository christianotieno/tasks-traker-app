package handlers

import (
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
	// TODO: Add dependencies, such as a task service or repository
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
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

	// TODO: Use the task service or repository to save the task

	w.WriteHeader(http.StatusCreated)
}
