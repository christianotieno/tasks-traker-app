package handlers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

// RouteHandler handles all the routes
func RouteHandler() {
	homeHandler := HomeHandler

	// Create a new router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/tasks", CreateTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/tasks", GetAllTasksHandler).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", GetTaskHandler).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", DeleteTaskHandler).Methods(http.MethodDelete)
	router.HandleFunc("/technicians", CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/technicians/{id}", GetUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/technician/{id}/tasks", GetAllTasksByUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler)

	// Redirect URLs with a trailing slash to the non-slash version
	router.PathPrefix("/tasks/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, "/"), http.StatusMovedPermanently)
	})

	// Start the server
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
