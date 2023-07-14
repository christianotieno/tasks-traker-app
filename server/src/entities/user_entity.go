package entities

type Role string

type User struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Tasks     *[]Task `json:"tasks"`
	ManagerID string  `json:"manager_id,omitempty"`
}

type UserJSON struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Tasks     *[]Task `json:"tasks"`
	ManagerID string  `json:"manager_id,omitempty"`
}
