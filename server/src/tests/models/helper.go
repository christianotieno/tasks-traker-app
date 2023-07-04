package models_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/christianotieno/tasks-traker-app/server/src/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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

func createTestUser(t *testing.T, db *sql.DB, firstName, lastName, email, role string) {
	t.Helper()
	_, err := db.Exec("INSERT INTO users (first_name, last_name, email, role) VALUES (?, ?, ?, ?)", firstName, lastName, email, role)
	assert.NoError(t, err)
}
