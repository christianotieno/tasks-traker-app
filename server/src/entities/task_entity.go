package entities

type Task struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
	Date    string `json:"date"`
	UserID  string `json:"user_id"`
}
