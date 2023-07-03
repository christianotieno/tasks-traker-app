package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	// Create a mock database connection
	db, _ := createMockDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Create a request with a task payload
	task := entities.Task{Summary: "Test Task", Date: "2023-07-01"}
	payload, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(payload))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Initialize the task model handler with the mock DB
	taskHandler := models.TaskHandler(db)

	// Call the CreateTask handler function
	taskHandler.CreateTask(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGetAllTasks(t *testing.T) {
	// Create a mock database connection
	db, _ := createMockDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Create a request
	req, _ := http.NewRequest("GET", "/tasks", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Initialize the task model handler with the mock DB
	taskHandler := models.TaskHandler(db)

	// Call the GetAllTasks handler function
	taskHandler.GetAllTasks(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var tasks []entities.Task
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		return
	}
	assert.Len(t, tasks, 2) // Assuming you have 2 tasks in the mock DB
}

func createMockDB() (*sql.DB, error) {
	// Open an in-memory SQLite database connection
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory database: %v", err)
	}

	return db, nil
}
