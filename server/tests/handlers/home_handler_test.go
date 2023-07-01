package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christianotieno/tasks-traker-app/server/handlers"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	handlers.HomeHandler(recorder, req)

	resp := recorder.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	expectedBody := "Hello, World!"
	body, err := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return
	}
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if string(body) != expectedBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedBody, string(body))
	}
}
