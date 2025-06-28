package payroll_periods

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "payroll/shared/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, log *logrus.Logger) {
    repo := NewPayrollPeriodRepository(db)
    service := NewPayrollPeriodService(repo)
    handler := NewPayrollPeriodHandler(service, log)

    routersGroup := router.Group("v1")
    {
        periodsGroup := routersGroup.Group("payroll_periods")
        periodsGroup.Use(utils.JWTAuthMiddleware())
        periodsGroup.POST("/create", handler.CreatePayrollPeriod)
    }
}