package models_tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

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

func createTestTask(t *testing.T, db *sql.DB, summary, date string) {
	t.Helper()
	_, err := db.Exec("INSERT INTO tasks (summary, date) VALUES (?, ?, ?)", summary, date)
	assert.NoError(t, err)
}

func createTestUser(t *testing.T) entities.User {
	t.Helper()
	db, err := config.TestDbConnect()
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Failed to close test database connection:", err)
		}
	}(db)

	// Generate random user data using gofakeit
	user := entities.User{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, false, false, 10),
		Role:      entities.Role(gofakeit.RandomString([]string{"Manager", "Technician"})),
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (first_name, last_name, email, password, role) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.Role)
	if err != nil {
		log.Fatal("Failed to create user in database:", err)
	}

	// Retrieve the created user from the database
	var lastInsertID int
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&lastInsertID)
	if err != nil {
		log.Fatal("Failed to retrieve LAST_INSERT_ID():", err)
	}

	// Retrieve the created user from the database
	err = db.QueryRow("SELECT * FROM users WHERE id = ?", lastInsertID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role)
	if err != nil {
		log.Fatal("Failed to retrieve created user from database:", err)
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

func generateUserData(num int) []entities.User {
	var users []entities.User
	pw := gofakeit.Password(true, true, true, false, false, 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < num; i++ {
		userData := entities.User{
			ID:        gofakeit.RandomInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
			Password:  string(hashedPassword),
			Role:      entities.Role(gofakeit.RandomString([]string{"Manager", "Technician"})),
		}

		users = append(users, userData)
	}

	return users
}

func generateTaskData(num int, userID int) []entities.Task {
	var tasks []entities.Task

	for i := 0; i < num; i++ {
		taskData := entities.Task{
			ID:      gofakeit.RandomInt([]int{1, 100}),
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
			ctx := context.WithValue(r.Context(), "userID", userData.ID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func generateValidToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	secret := []byte(os.Getenv("SECRET"))
	tokenString, _ := token.SignedString(secret)

	return tokenString
}
