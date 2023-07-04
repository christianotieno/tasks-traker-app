package models_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/christianotieno/tasks-traker-app/server/src/config"
	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskModel(t *testing.T) {
	db, err := config.TestDbConnect()
	if err != nil {
		t.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(db)

	taskModel := models.TaskModel{Db: db}

	t.Run("CreateTask", func(t *testing.T) {
		reqBody := []byte(`{"summary": "Test Summary", "date": "2023-01-01"}`)
		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		taskModel.CreateTask(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdTask entities.Task
		err = json.Unmarshal(w.Body.Bytes(), &createdTask)
		assert.NoError(t, err)

		assert.NotNil(t, createdTask.ID)
		assert.Equal(t, "Test Summary", createdTask.Summary)
		assert.Equal(t, "2023-01-01 00:00:00", createdTask.Date)
	})

	t.Run("GetAllTasks", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		taskModel.GetAllTasks(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var tasks []entities.Task
		err = json.Unmarshal(w.Body.Bytes(), &tasks)
		assert.NoError(t, err)

		assert.NotEmpty(t, tasks)
	})

	t.Run("GetTask", func(t *testing.T) {
		_, err := db.Exec("DELETE FROM tasks WHERE id = ?", 1)
		assert.NoError(t, err)

		_, err = db.Exec("INSERT INTO tasks (id, summary, date) VALUES (?, ?, ?)", 1, "Test Summary", "2023-01-01")
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/tasks/1", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		taskModel.GetTask(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var task entities.Task
		err = json.Unmarshal(w.Body.Bytes(), &task)
		assert.NoError(t, err)

		assert.Equal(t, "Test Summary", task.Summary)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		_, err := db.Exec("DELETE FROM tasks WHERE id = ?", 1)
		assert.NoError(t, err)

		_, err = db.Exec("INSERT INTO tasks (id, summary, date) VALUES (?, ?, ?)", 1, "Test Summary", "2023-01-01")
		assert.NoError(t, err)

		reqBody := []byte(`{"summary": "Updated Summary", "date": "2023-02-02"}`)
		req, err := http.NewRequest(http.MethodPatch, "/tasks/1", bytes.NewReader(reqBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		taskModel.UpdateTask(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var summary string
		var date string
		err = db.QueryRow("SELECT summary, date FROM tasks WHERE id = ?", 1).Scan(&summary, &date)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Summary", summary)
		assert.Equal(t, "2023-02-02 00:00:00", date)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		taskModel.DeleteTask(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Message string `json:"message"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Task deleted successfully", response.Message)
	})
}
