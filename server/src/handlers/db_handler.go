package handlers

import (
	"database/sql"
	"github.com/christianotieno/tasks-traker-app/server/src/config"
	"log"
)

// openDbConnection connects the database and returns a database connection
func openDbConnection() (*sql.DB, error) {
	// Connect to the database
	db, err := config.DbConnect()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

// closeDbConnection closes the database connection
func closeDbConnection(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
