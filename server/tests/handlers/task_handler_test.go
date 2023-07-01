package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/christianotieno/tasks-traker-app/server/handlers"
)

func TestCreateTask(t *testing.T) {
	// Given
	taskHandler := handlers.NewTaskHandler()

	// Task payload
	task := handlers.Task{
		Summary: "Perform maintenance",
		Date:    time.Now(),
	}
	payload, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task payload: %v", err)
	}

	// New request with the task payload
	req, err := http.NewRequest("POST", "/tasks", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Response recorder
	recorder := httptest.NewRecorder()

	// When
	taskHandler.CreateTask(recorder, req)

	// Then
	resp := recorder.Result()

	// Check response status code
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	// TODO: Add additional assertions or checks for more future logic
}
