package models_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christianotieno/tasks-traker-app/server/src/models"
)

func TestCreateUser(t *testing.T) {
	t.Run("InvalidRequestMethod", func(t *testing.T) {
		// Given
		db := setupTestDB(t)
		um := &models.UserModel{Db: db}
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		// When
		um.CreateUser(rr, req)

		// Then
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})
}
