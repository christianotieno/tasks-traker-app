package entities

type Task struct {
	ID      int64  `json:"id"`
	Summary string `json:"summary"`
	Date    string `json:"date"`
}
