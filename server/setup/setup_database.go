package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase() {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/task_manager")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database connection:", err)
		}
	}(db)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS task_manager")
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	sqlFile, err := ioutil.ReadFile("./setup/setup.sql")
	if err != nil {
		log.Fatal("Failed to read SQL file:", err)
	}

	// Split the SQL commands by semicolon
	sqlCommands := strings.Split(string(sqlFile), ";")

	// Execute each SQL command
	for _, cmd := range sqlCommands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		_, err := db.Exec(cmd)
		if err != nil {
			log.Fatal("Failed to execute SQL command:", err)
		}
	}

	fmt.Println("Database tables created successfully!")
}

func main() {
	setupDatabase()
}
