package main

import (
	"log"

	"github.com/christianotieno/tasks-traker-app/server/src/handlers"
)

func main() {
	err := handlers.InitDbConnection()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		err := handlers.CloseDbConnection()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	handlers.RouteHandler()
}
