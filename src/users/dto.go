package users


type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        int  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"not null" json:"username"`
	Token string    	`json:"token"`
}
