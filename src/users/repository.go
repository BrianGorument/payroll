package users

import (
	"gorm.io/gorm"
)

// userRepository struct
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository (Dependency Injection)
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

// Implementasi UserRepository
func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByUUID(id int) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}