package handlers

import (
	"github.com/christianotieno/tasks-traker-app/server/src/models"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userHandler := models.UserHandler(db)
	userHandler.CreateUser(w, r)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	userHandler := models.UserHandler(db)
	userHandler.GetUser(w, r, userID)
}

func GetAllTasksByUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	technicianID := vars["id"]
	userHandler := models.UserHandler(db)
	userHandler.GetAllTasksByUserID(w, r, technicianID)
}

func GetAllUsersAndAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	managerID := vars["id"]
	userHandler := models.UserHandler(db)
	userHandler.GetAllUsersAndAllTasks(w, r, managerID)
}
