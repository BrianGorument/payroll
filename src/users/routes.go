package users

import (
	"payroll/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RegisterRoutes
func RegisterRoutes(router *gin.Engine, db *gorm.DB, log *logrus.Logger) {
	repo := NewUserRepository(db)
	service := NewUserService(repo)
	handler := NewUserHandler(service, log)

	routersGroup := router.Group("v1")
	{
		usersGroup := routersGroup.Group("users")

		usersGroup.GET("/", utils.JWTAuthMiddleware(), handler.GetAllUsers)
		usersGroup.POST("/login", handler.LoginUser)
	}
}
