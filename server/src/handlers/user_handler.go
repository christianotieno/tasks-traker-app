package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/christianotieno/tasks-traker-app/server/src/models"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	userHandler := models.UserHandler(db)
	userHandler.CreateUser(w, r)
}

func GetAllTasksByUserHandler(w http.ResponseWriter, r *http.Request) {
	userHandler := models.UserHandler(db)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userHandler.GetAllTasksByUserID(w, r, mux.Vars(r)["id"])
	})
	authenticate(handler).ServeHTTP(w, r)
}

func GetAllUsersAndAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	userHandler := models.UserHandler(db)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userHandler.GetAllUsersAndAllTasks(w, r)
	})
	authenticate(handler).ServeHTTP(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	taskHandler := models.UserHandler(db)
	taskHandler.Login(w, r)
}
