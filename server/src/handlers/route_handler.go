package handlers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// RouteHandler handles all the routes
func RouteHandler() {
	homeHandler := HomeHandler

	// Create a new router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/tasks", CreateTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/tasks", GetAllTasksHandler).Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler)

	log.Println("Server listening on http://localhost:8080")
	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
