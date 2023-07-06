package models_tests

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/gorilla/mux"
)

func TestCreateUser(t *testing.T) {
	secret := os.Getenv("SECRET")
	userData := entities.User{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, false, false, 10),
		Role:      entities.Role(gofakeit.RandomString([]string{"Manager", "Technician"})),
	}

	err := os.Setenv("SECRET", secret)
	if err != nil {
		return
	}

	t.Run("Should create user successfully ", func(t *testing.T) {
		db := setupTestDB(t)
		userModel := models.UserModel{Db: db}

		// Given
		user, err := stringifyUser(userData)
		if err != nil {
			log.Fatalf("Failed to stringify user: %v", err)
		}

		// When
		req, err := http.NewRequest(http.MethodPost, "/users", strings.NewReader(user))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := httptest.NewRecorder()

		userModel.CreateUser(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
		}

		expectedUser := entities.User{
			ID:        userData.ID,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Email:     userData.Email,
			Role:      userData.Role,
		}

		responseBody := w.Body.String()

		var actualResponse struct {
			User  entities.User `json:"user"`
			Token string        `json:"token"`
		}

		err = json.Unmarshal([]byte(responseBody), &actualResponse)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expectedUserCopy := expectedUser
		expectedUserCopy.ID = actualResponse.User.ID

		if !reflect.DeepEqual(actualResponse.User, expectedUserCopy) {
			t.Errorf("Expected user %+v, but got %+v", expectedUserCopy, actualResponse.User)
		}

		cleanupTestDB(t, db)
	})

	t.Run("Should return an error when user already exists", func(t *testing.T) {
		db := setupTestDB(t)
		userModel := models.UserModel{Db: db}

		_, err := db.Exec("INSERT INTO users (first_name, last_name, email, password, role) VALUES (?, ?, ?, ?, ?)",
			userData.FirstName, userData.LastName, userData.Email, "password", userData.Role)
		if err != nil {
			t.Fatalf("Failed to insert user into database: %v", err)
		}

		// Given
		user, err := stringifyUser(userData)
		if err != nil {
			t.Fatalf("Failed to stringify user: %v", err)
		}

		// When
		req, err := http.NewRequest(http.MethodPost, "/users", strings.NewReader(user))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := httptest.NewRecorder()

		userModel.CreateUser(w, req)

		// Then
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
		}

		expectedErrorMessage := "Email already exists, please try again with a different email"
		responseBody := w.Body.String()
		if !strings.Contains(responseBody, expectedErrorMessage) {
			t.Errorf("Expected error message '%s' in response body, but got: %s", expectedErrorMessage, responseBody)
		}
		cleanupTestDB(t, db)
	})
}

func TestGetAllTasksByUserID(t *testing.T) {
	db := setupTestDB(t)

	userModel := models.UserModel{Db: db}
	userData := generateUserData(1)[0]
	taskData := generateTaskData(5, userData.ID)

	req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()

	// Set up the user data
	_, err = db.Exec("INSERT INTO users (id, first_name, last_name, email, password, role) VALUES (?, ?, ?, ?, ?, ?)",
		userData.ID, userData.FirstName, userData.LastName, userData.Email, userData.Password, userData.Role)
	if err != nil {
		t.Fatalf("Failed to set up user data: %v", err)
	}

	// Set up the user task data
	for _, task := range taskData {
		_, err = db.Exec("INSERT INTO tasks ( user_id, summary, date) VALUES (?,?,?)",
			task.UserID, task.Summary, task.Date)
	}

	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	// Set up the user ID in the request context
	ctx := context.WithValue(req.Context(), "userID", userData.ID)
	req = req.WithContext(ctx)

	// Generate a valid JWT token
	token := generateValidToken(userData.ID)

	req.Header.Set("Authorization", token)

	router := mux.NewRouter()

	router.Use(MockAuthorizationMiddleware(userData))

	router.HandleFunc("/tasks", userModel.GetAllTasksByUserID)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var responseTasks []entities.Task
	err = json.Unmarshal(w.Body.Bytes(), &responseTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if len(responseTasks) == 0 {
		t.Error("Expected tasks, but got an empty response")
	}
	cleanupTestDB(t, db)
}
