package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/christianotieno/tasks-traker-app/server/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func TestCreateTask(t *testing.T) {
	// Establish a database connection
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create a new task handler with the database connection
	taskHandler := handlers.NewTaskHandler(db)

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

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/task_manager_test")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func teardownTestDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		return
	}
}
