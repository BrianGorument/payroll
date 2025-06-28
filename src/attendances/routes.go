package attendances

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "payroll/shared/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, log *logrus.Logger) {
    repo := NewAttendanceRepository(db)
    service := NewAttendanceService(repo)
    handler := NewAttendancesHandler(service, log)

    routersGroup := router.Group("v1")
    {
        attendanceGroup := routersGroup.Group("attendances")
        attendanceGroup.Use(utils.JWTAuthMiddleware())
        attendanceGroup.POST("/create", handler.CreateAttendances)
    }
}