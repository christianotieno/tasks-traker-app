package handlers

import (
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"log"
	"net/http"
)

// CreateTaskHandler defines the route handler function for creating a task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := openDbConnection()
	if err != nil {
		log.Fatal(err)
		return
	}

	taskHandler := models.TaskHandler(db)
	taskHandler.CreateTask(w, r)

	err = closeDbConnection(db)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// GetAllTasksHandler defines the route handler function for retrieving all tasks
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDbConnection()
	if err != nil {
		log.Fatal(err)
		return
	}

	taskHandler := models.TaskHandler(db)
	taskHandler.GetAllTasks(w, r)

	err = closeDbConnection(db)
	if err != nil {
		log.Fatal(err)
		return
	}
}
