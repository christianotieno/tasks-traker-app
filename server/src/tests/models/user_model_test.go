package models_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	userModel := models.UserModel{Db: db}

	t.Run("CreateUser", func(t *testing.T) {
		reqBody := []byte(`{"first_name": "John", "last_name": "Doe", "email": "john.doe@mail.com", "role": "Manager"}`)
		req := createTestRequest(t, http.MethodPost, "/users", reqBody)

		w := httptest.NewRecorder()

		userModel.CreateUser(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdUser entities.User
		unmarshalResponse(t, w.Body.Bytes(), &createdUser)

		assert.NotNil(t, createdUser.ID)
		assert.Equal(t, "John Doe", createdUser.FirstName+" "+createdUser.LastName)
		assert.Equal(t, "john.doe@mail.com", createdUser.Email)
		assert.Equal(t, "Manager", createdUser.Role)
	})

	t.Run("GetUser", func(t *testing.T) {
		createTestUser(t, db, "Jane", "Doe", "jane.doe@mail.com", "Technician")

		req := createTestRequest(t, http.MethodGet, "/users/1", nil)

		w := httptest.NewRecorder()

		userModel.GetUser(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var user entities.User
		unmarshalResponse(t, w.Body.Bytes(), &user)

		assert.Equal(t, "Jane Doe", user.FirstName+" "+user.LastName)
		assert.Equal(t, "jane.doe@mail.com", user.Email)
		assert.Equal(t, "Technician", user.Role)
	})

	t.Run("GetAllTasksByUserID", func(t *testing.T) {
		createTestUser(t, db, "Jane", "Doe", "jane.doe@mail.com", "Technician")
		createTestTask(t, db, "Performed Task 1", "2023-06-04")
		createTestTask(t, db, "Performed Task 2", "2023-07-04")

		req := createTestRequest(t, http.MethodGet, "/users/{id}/tasks", nil)

		w := httptest.NewRecorder()

		userModel.GetAllTasksByUserID(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var user []entities.User
		unmarshalResponse(t, w.Body.Bytes(), &user)

		assert.NotEmpty(t, user)
		assert.Equal(t, "Performed Task 1", user[0].Tasks[0].Summary)
		assert.Equal(t, "Performed Task 2", user[0].Tasks[1].Summary)
	})

	t.Run("GetAllUsersAndAllTasks", func(t *testing.T) {
		createTestUser(t, db, "Mathew", "West", "mathew.west@mail.com", "Manager")
		createTestUser(t, db, "Peter", "Smith", "peter.smith@mail.com", "Technician")
		createTestUser(t, db, "Jane", "Doe", "jane.doe@mail.com", "Technician")
		createTestTask(t, db, "Performed Task 1", "2023-06-04")
		createTestTask(t, db, "Performed Task 2", "2023-07-04")
		createTestTask(t, db, "Performed Task 3", "2023-06-04")

		req := createTestRequest(t, http.MethodGet, "/managers/{id}/users", nil)

		w := httptest.NewRecorder()

		userModel.GetAllTasksByUserID(w, req, "1")

		assert.Equal(t, http.StatusOK, w.Code)

		var users []entities.User
		unmarshalResponse(t, w.Body.Bytes(), &users)

		assert.NotEmpty(t, users)
		assert.Equal(t, "Performed Task 1", users[0].Tasks[0].Summary)
		assert.Equal(t, "Performed Task 2", users[0].Tasks[1].Summary)
		assert.Equal(t, "Performed Task 3", users[1].Tasks[0].Summary)
	})
}
