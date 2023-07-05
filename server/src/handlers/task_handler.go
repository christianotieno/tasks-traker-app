package handlers

import (
	"net/http"

	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/gorilla/mux"
)

// CreateTaskHandler defines the route handler function for creating a task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskHandler := models.TaskHandler(db)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskHandler.CreateTask(w, r)
	})
	authenticate(handler).ServeHTTP(w, r)
}

// DeleteTaskHandler defines the route handler function for deleting a task
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	taskHandler := models.TaskHandler(db)
	taskHandler.DeleteTask(w, r, id)
}
