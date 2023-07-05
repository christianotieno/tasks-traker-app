package entities

type Role string

const (
	Technician Role = "Technician"
	Manager    Role = "Manager"
)

type User struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Role      Role    `json:"role"`
	Tasks     *[]Task `json:"tasks"`
	Password  string  `json:"-"`
}
