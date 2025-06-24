package users

import (
	"errors"
	"fmt"
	"payroll/shared/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

// userService struct
type userService struct {
	repo IUserRepository
}

// NewUserService (Dependency Injection)
func NewUserService(repo IUserRepository) IUserService {
	return &userService{repo}
}

// Implementasi UserService
func (s *userService) RegisterUser(req CreateUserRequest) (*UserResponse, error) {
	existingUsers, _ := s.repo.FindAll()
	for _, u := range existingUsers {
		if u.Email == req.Email {
			return nil, errors.New("email already registered")
		}
	}
	hashedPassword, _ := utils.HashPassword(req.Password)
	validate := validator.New()
	err := validate.Var(req.Email, "required,email")
	if err != nil {
		fmt.Println("Invalid email:", err)
		return nil, errors.New("invalid Email Format")
	}

	timenow := time.Now()

	user := User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role_code: req.Role_code,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: timenow,
		UpdatedAt: timenow,
	}
	err = s.repo.Create(&user)
	if err != nil {
		return nil, err
	}

	return &UserResponse{ID: user.ID, Username: user.Username, Email: user.Email}, nil
}

func (s *userService) GetAllUsers() ([]UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return userResponses, nil
}

func (s *userService) LoginUser(req CreateUserRequest) (*UserResponse, error) {
	existingUsers, _ := s.repo.FindAll()
	for _, u := range existingUsers {
		if u.Email == req.Email {
			return &UserResponse{
				ID:       u.ID,
				Username: u.Username,
				Email:    u.Email,
			}, nil
		}
	}

	return nil, errors.New("user not found")
}
