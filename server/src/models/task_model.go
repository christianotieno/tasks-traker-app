package models

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
)

type TaskModel struct {
	Db *sql.DB
}

func TaskHandler(db *sql.DB) *TaskModel {
	return &TaskModel{
		Db: db,
	}
}

func (tm *TaskModel) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		return
	}

	var task entities.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	// Insert the task into the database
	_, err = tm.Db.Exec("INSERT INTO tasks (summary, date) VALUES (?, ?)", task.Summary, task.Date)
	if err != nil {
		http.Error(w, "Task creation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (tm *TaskModel) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := tm.Db.Query("SELECT * FROM tasks")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}(rows)

	var tasks []entities.Task

	// Retrieve tasks from the database
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(&task.ID, &task.Summary, &task.Date)
		if err != nil {
			http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		tasks = append(tasks, task)
	}

	// Serialize tasks to JSON
	response, err := json.Marshal(&tasks)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
}
