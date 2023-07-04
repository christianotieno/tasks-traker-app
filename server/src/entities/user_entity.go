package entities

type Role string

const (
	Technician Role = "Technician"
	Manager    Role = "Manager"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	Tasks     []Task `json:"tasks"`
}