package main

import (
	"github.com/christianotieno/tasks-traker-app/server/src/handlers"
	"log"
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
