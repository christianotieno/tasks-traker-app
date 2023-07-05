package entities

type Task struct {
	ID         int    `json:"id"`
	Summary    string `json:"summary"`
	Date       string `json:"date"`
	UserID     int    `json:"user_id"`
	Manager    *User  `json:"manager,omitempty"`
	Technician *User  `json:"technician,omitempty"`
}
