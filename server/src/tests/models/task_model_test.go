package models_tests

import (
	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskModel(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	taskModel := models.TaskModel{Db: db}

	t.Run("CreateTask", func(t *testing.T) {
		reqBody := []byte(`{"summary": "Test Summary", "date": "2023-01-01"}`)
		req := createTestRequest(t, http.MethodPost, "/tasks", reqBody)

		w := httptest.NewRecorder()

		taskModel.CreateTask(w, req, "1")

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdTask entities.Task
		unmarshalResponse(t, w.Body.Bytes(), &createdTask)

		assert.NotNil(t, createdTask.ID)
		assert.Equal(t, "Test Summary", createdTask.Summary)
		assert.Equal(t, "2023-01-01 00:00:00", createdTask.Date)
	})

	t.Run("GetAllTasks", func(t *testing.T) {
		req := createTestRequest(t, http.MethodGet, "/tasks", nil)

		w := httptest.NewRecorder()

		taskModel.GetAllTasks(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var tasks []entities.Task
		unmarshalResponse(t, w.Body.Bytes(), &tasks)

		assert.NotEmpty(t, tasks)
	})

	t.Run("GetTask", func(t *testing.T) {
		createTestTask(t, db, "Test Summary", "2023-01-01")

		req := createTestRequest(t, http.MethodGet, "/tasks/1", nil)

		w := httptest.NewRecorder()

		taskModel.GetTask(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var task entities.Task
		unmarshalResponse(t, w.Body.Bytes(), &task)

		assert.Equal(t, "Test Summary", task.Summary)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		createTestTask(t, db, "Test Summary", "2023-01-01")

		reqBody := []byte(`{"summary": "Updated Summary", "date": "2023-02-02"}`)
		req := createTestRequest(t, http.MethodPatch, "/tasks/1", reqBody)

		w := httptest.NewRecorder()

		taskModel.UpdateTask(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var summary string
		var date string
		err := db.QueryRow("SELECT summary, date FROM tasks WHERE id = ?", 1).Scan(&summary, &date)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Summary", summary)
		assert.Equal(t, "2023-02-02 00:00:00", date)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		createTestTask(t, db, "Test Summary", "2023-01-01")

		req := createTestRequest(t, http.MethodDelete, "/tasks/1", nil)

		w := httptest.NewRecorder()

		taskModel.DeleteTask(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Message string `json:"message"`
		}
		unmarshalResponse(t, w.Body.Bytes(), &response)

		assert.Equal(t, "Task deleted successfully", response.Message)
	})
}
