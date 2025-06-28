package overtimes

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "payroll/shared/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, log *logrus.Logger) {
    repo := NewOvertimeRepository(db)
    service := NewOvertimeService(repo)
    handler := NewOvertimeHandler(service, log)

    routersGroup := router.Group("v1")
    {
        overtimeGroup := routersGroup.Group("overtimes")
        overtimeGroup.Use(utils.JWTAuthMiddleware())
        overtimeGroup.POST("/create", handler.CreateOvertime)
    }
}