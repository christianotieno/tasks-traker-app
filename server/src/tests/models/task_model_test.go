package models_tests

import (
	"context"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTask(t *testing.T) {
	t.Run("InvalidRequestMethod", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		req, err := http.NewRequest("GET", "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		tm.CreateTask(rr, req)

		// Then
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("MissingUserID", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		req, err := http.NewRequest("POST", "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		tm.CreateTask(rr, req)

		// Then
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("Test: foreign key constraint", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		reqBody := `{"summary": "Sample Task", "date": "2023-07-06"}`
		req, err := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		ctx := context.WithValue(req.Context(), "userID", "123")
		ctx = context.WithValue(ctx, "managerID", "456")
		req = req.WithContext(ctx)

		// When
		tm.CreateTask(rr, req)

		// Then
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		reqBody := `{"summary": "Sample Task", "date": "2023-07-06"}`
		req, err := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		ctx := context.WithValue(req.Context(), "userID", "123")
		req = req.WithContext(ctx)

		// When
		tm.CreateTask(rr, req)

		// Then
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, rr.Code)
		}
	})

}

func TestDeleteTask(t *testing.T) {
	t.Run("InvalidRequestMethod", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		req, err := http.NewRequest("GET", "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		tm.DeleteTask(rr, req, "123")

		// Then
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("MissingUserID", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		req, err := http.NewRequest("DELETE", "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		tm.DeleteTask(rr, req, "123")

		// Then
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
		}
	})

}

func TestUpdateTask(t *testing.T) {
	t.Run("InvalidRequestMethod", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		req, err := http.NewRequest("GET", "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		tm.UpdateTask(rr, req, "123")

		// Then
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("MissingUserID", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		tm := &models.TaskModel{Db: db}
		req, err := http.NewRequest("PUT", "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		tm.UpdateTask(rr, req, "123")

		// Then
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

}
