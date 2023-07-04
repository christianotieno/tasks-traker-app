package handlers

import (
	"database/sql"
	"github.com/christianotieno/tasks-traker-app/server/src/config"
)

var db *sql.DB // Declare a global variable for the database connection

// InitDbConnection initializes the database connection
func InitDbConnection() error {
	var err error
	db, err = config.DbConnect()
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// CloseDbConnection closes the database connection
func CloseDbConnection() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
