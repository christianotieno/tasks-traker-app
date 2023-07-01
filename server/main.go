package main

import (
	"database/sql"
	"fmt"
	"github.com/christianotieno/tasks-traker-app/server/handlers"
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

	taskHandler := handlers.NewTaskHandler()
	homeHandler := handlers.HomeHandler

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", taskHandler.CreateTask)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
