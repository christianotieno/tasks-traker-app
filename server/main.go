package main

import (
	"database/sql"
	"fmt"
	"github.com/christianotieno/tasks-traker-app/server/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/tasks_tracker")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	taskHandler := handlers.NewTaskHandler(db)
	homeHandler := handlers.HomeHandler

	// Create a new router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", taskHandler.ListTasks).Methods(http.MethodGet)
	router.HandleFunc("/", homeHandler)

	fmt.Println("Server listening on http://localhost:8080")
	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
