package models_tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"

	"github.com/christianotieno/tasks-traker-app/server/src/config"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := config.TestDbConnect()
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	tables := []string{
		"tasks",
		"users",
		"managers",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Fatalf("Failed to delete data from table %s: %v", table, err)
		}
	}

	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func createTestRequest(t *testing.T, method, url string, body []byte) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	assert.NoError(t, err)
	return req
}

func unmarshalResponse(t *testing.T, data []byte, v interface{}) {
	t.Helper()
	err := json.Unmarshal(data, v)
	assert.NoError(t, err)
}

func createTestTask(t *testing.T, db *sql.DB, userID string) entities.Task {
	t.Helper()

	task := entities.Task{
		Summary: "Test Task",
		Date:    "2023-07-05",
		UserID:  userID,
	}

	_, err := db.Exec("INSERT INTO tasks (summary, date, user_id) VALUES (?, ?, ?)",
		task.Summary, task.Date, task.UserID)
	if err != nil {
		t.Fatal("Failed to create task in database:", err)
	}

	err = db.QueryRow("SELECT * FROM tasks WHERE summary = ?", task.Summary).Scan(
		&task.ID, &task.Summary, &task.Date, &task.UserID)
	if err != nil {
		t.Fatal("Failed to retrieve created task from database:", err)
	}

	return task
}

func createTestUser(t *testing.T) entities.UserJSON {
	t.Helper()

	// Generate random user data using gofakeit
	user := entities.UserJSON{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, false, false, 10),
	}

	id := uuid.New().String()

	db, err := config.TestDbConnect()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Log("Failed to close test database connection:", err)
		}
	}()

	_, err = db.Exec("INSERT INTO users (id, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)",
		id, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		t.Fatal("Failed to create user in database:", err)
	}

	err = db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		t.Fatal("Failed to retrieve created user from database:", err)
	}

	return user
}

func stringifyUser(user entities.User) (string, error) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("Failed to stringify user:", err)
		return "", err
	}
	return string(userJSON), nil
}

func generateUserData(num int) []entities.UserJSON {
	var users []entities.UserJSON
	pw := gofakeit.Password(true, true, true, false, false, 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < num; i++ {
		userData := entities.UserJSON{
			ID:        uuid.New().String(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
			Password:  string(hashedPassword),
		}

		users = append(users, userData)
	}

	return users
}

func generateTaskData(num int, userID string) []entities.Task {
	var tasks []entities.Task
	for i := 0; i < num; i++ {
		taskData := entities.Task{
			ID:      uuid.New().String(),
			Summary: gofakeit.Sentence(5),
			Date:    gofakeit.Date().Format("2006-01-02"),
			UserID:  userID,
		}

		tasks = append(tasks, taskData)
	}

	return tasks
}

func MockAuthorizationMiddleware(userData entities.User) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "userID", string(userData.ID))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func generateValidToken(userID []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": string(userID),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	secret := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Println("Failed to generate token:", err)
		return "", err
	}

	return tokenString, nil
}
