package users


type IUserService interface {
	//	RegisterUser(req UserLoginRequest) (*UserResponse, error)
	LoginUser(req UserLoginRequest) (*UserResponse, error)
	GetAllUsers() ([]UserResponse, error)
}

type IUserRepository interface {
	//	Create(user *User) error
	FindAll() ([]User, error)
	FindByUUID(id int) (*User, error)
	FindByUsername(username string) (*User, error)
}
