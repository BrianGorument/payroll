package users

import (
	"errors"
	"fmt"
	"payroll/shared/utils"
)

// userService struct
type userService struct {
	repo IUserRepository
}

// NewUserService (Dependency Injection)
func NewUserService(repo IUserRepository) IUserService {
	return &userService{repo}
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
		})
	}

	return userResponses, nil
}

func (s *userService) LoginUser(req UserLoginRequest) (*UserResponse, error) {
	eu, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}	
	verifiedPassword :=  utils.VerifyPassword(req.Password, eu.Password,)
	
	if !verifiedPassword {
		return nil, errors.New("invalid password")
	}
	
	token, err := utils.CreateJWTToken(eu.ID, eu.Username , eu.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}
	
	existingUsers, _ := s.repo.FindAll()
	for _, u := range existingUsers {
		if u.Username == req.Username {
			return &UserResponse{
				ID:       u.ID,
				Username: u.Username,
				Token: token,
			},  nil
		}
	}

	return nil, errors.New("user not found")
}
