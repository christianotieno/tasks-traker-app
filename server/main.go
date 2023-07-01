package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Task struct {
	Summary string    `json:"summary"`
	Date    time.Time `json:"date"`
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", createTaskHandler)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var task Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// TODO: Save the task to MySQL

	w.WriteHeader(http.StatusCreated)
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		return
	}
}
