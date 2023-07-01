package main

import (
	"fmt"
	"github.com/christianotieno/tasks-traker-app/server/handlers"
	"log"
	"net/http"
)

func main() {
	taskHandler := handlers.NewTaskHandler()
	homeHandler := handlers.HomeHandler

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", taskHandler.CreateTask)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
