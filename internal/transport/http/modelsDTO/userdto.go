package modelsDTO

type UsersDTO struct {
	Users []UserDTO `json:"users"`
	Total int       `json:"total"`
}

type UserDTO struct {
	ID       int64    `json:"ID"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	LastName string   `json:"last_name"`
	Role     []string `json:"role"`
}
type CreateUserDTO struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	LastName string   `json:"last_name"`
	Role     []string `json:"role"`
}
type UpdateUserDTO struct {
	ID       int64    `json:"ID"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	LastName string   `json:"last_name"`
	Role     []string `json:"role"`
}
