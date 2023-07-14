package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/christianotieno/tasks-traker-app/server/src/handlers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = handlers.InitDbConnection()
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

	// Start the server
	go func() {
		handlers.RouteHandler()
	}()

	// Start Kafka consumer
	go handlers.HandleKafkaMessages([]string{"localhost:9092"})

	// Keep the main function running
	select {}
}
