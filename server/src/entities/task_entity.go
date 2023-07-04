package entities

type Task struct {
	ID         int64  `json:"id"`
	Summary    string `json:"summary"`
	Date       string `json:"date"`
	UserID     int64  `json:"user_id"`
	Manager    *User  `json:"manager,omitempty"`
	Technician *User  `json:"technician,omitempty"`
}
