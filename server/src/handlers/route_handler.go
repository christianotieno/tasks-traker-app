package handlers

import (
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// RouteHandler handles all the routes
func RouteHandler() {
	homeHandler := HomeHandler

	// Create a new router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/login", LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/users", CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/tasks", CreateTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", UpdateTaskHandler).Methods(http.MethodPatch)
	router.HandleFunc("/tasks/{id}", DeleteTaskHandler).Methods(http.MethodDelete)
	router.HandleFunc("/users", GetAllUsersAndAllTasksHandler).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}/tasks", GetAllTasksByUserHandler).Methods(http.MethodGet)

	// Redirect URLs with a trailing slash to the non-slash version
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, "/"), http.StatusMovedPermanently)
	})))

	// Start the server
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
