package models_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	// Given
	db := setupTestDB(t)

	createTestUser(t)
	tm := &models.TaskModel{Db: db}
	// Create a new HTTP request
	reqBody := `{
		"summary": "Sample Task",
		"date": "2023-07-06"
	}`
	req, err := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the userID in the request context
	ctx := context.WithValue(req.Context(), "userID", float64(123))

	// Assign the context with userID to the request
	req = req.WithContext(ctx)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the CreateTask function
	tm.CreateTask(rr, req)

	// Check the response status code
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, rr.Code)
	}

	// Verify the response body
	expectedResponseBody := `{"id": 1, "userID": 123, "summary": "Sample Task", "date": "2023-07-06"}`
	if rr.Body.String() != expectedResponseBody {
		t.Errorf("Expected response body %s, but got %s", expectedResponseBody, rr.Body.String())
	}

	// Additional assertions can be added to check the behavior of the function
}

func TestDeleteTask(t *testing.T) {
	// Given
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)
	taskModel := &models.TaskModel{Db: db}

	task := entities.Task{
		Summary: "Test Task",
		Date:    "2023-07-05",
		UserID:  1,
	}

	res, err := taskModel.Db.Exec("INSERT INTO tasks (user_id, summary, date) VALUES (?, ?, ?)",
		task.UserID, task.Summary, task.Date)
	if err != nil {
		t.Fatalf("Failed to insert test task into the database: %v", err)
	}

	taskID, getErr := res.LastInsertId()

	if getErr != nil {
		t.Fatalf("Failed to retrieve task ID: %v", getErr)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", taskID), nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}

	w := httptest.NewRecorder()

	// When
	taskModel.DeleteTask(w, req, strconv.FormatInt(taskID, 10))

	// Then
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Message string `json:"message"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check the response message
	expectedMessage := "Task deleted successfully"
	if response.Message != expectedMessage {
		t.Errorf("Expected message '%s', but got '%s'", expectedMessage, response.Message)
	}

	var count int
	err = taskModel.Db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ?", taskID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	assert.Equal(t, 0, count)
}
