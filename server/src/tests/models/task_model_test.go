package models_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	// Given
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create a new task model instance
	taskModel := &models.TaskModel{Db: db}

	reqBody := []byte(`{"summary": "Test Summary", "date": "2023-01-01"}`)

	req := createTestRequest(t, http.MethodPost, "/tasks", reqBody)

	// When
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "userID", "1")

	req = req.WithContext(ctx)

	taskModel.CreateTask(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTask entities.Task

	err := json.Unmarshal(w.Body.Bytes(), &createdTask)

	assert.NoError(t, err)

	// Then
	assert.NotNil(t, createdTask.ID)
	assert.Equal(t, "Test Summary", createdTask.Summary)
	assert.Equal(t, "2023-01-01", createdTask.Date)
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
