package handlers

import (
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateTaskHandler defines the route handler function for creating a task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	technicianID := vars["id"]
	taskHandler := models.TaskHandler(db)
	taskHandler.CreateTask(w, r, technicianID)
}

// GetAllTasksHandler defines the route handler function for retrieving all tasks
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	taskHandler := models.TaskHandler(db)
	taskHandler.GetAllTasks(w, r)
}

// GetTaskHandler defines the route handler function for retrieving a task
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	taskHandler := models.TaskHandler(db)
	taskHandler.GetTask(w, r, id)
}

// UpdateTaskHandler defines the route handler function for updating a task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	taskHandler := models.TaskHandler(db)
	taskHandler.UpdateTask(w, r, id)
}

// DeleteTaskHandler defines the route handler function for deleting a task
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	taskHandler := models.TaskHandler(db)
	taskHandler.DeleteTask(w, r, id)
}
