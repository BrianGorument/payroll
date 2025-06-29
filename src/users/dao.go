package users

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`   
	Username  string    `gorm:"not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Salary    float64   `gorm:"type:decimal(15,2)" json:"salary"`
	Role      string    `gorm:"type:user_role;not null" json:"role"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy string    `gorm:"not null" json:"created_by"`
	UpdatedBy string    `gorm:"not null" json:"updated_by"`
}

func (User) TableName() string {
	return "users"
}
