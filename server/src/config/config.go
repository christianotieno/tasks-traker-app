package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func DbConnect() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/task_manager")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestDbConnect() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/task_manager_test")
	if err != nil {
		return nil, err
	}
	return db, nil
}
