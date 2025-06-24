package users

type IUserService interface {
	RegisterUser(req CreateUserRequest) (*UserResponse, error)
	LoginUser(req CreateUserRequest) (*UserResponse, error)
	GetAllUsers() ([]UserResponse, error)
}

type IUserRepository interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
}
