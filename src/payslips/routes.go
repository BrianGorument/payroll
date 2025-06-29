package payslips

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "payroll/shared/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, log *logrus.Logger) {
    repo := NewPayslipRepository(db)
    service := NewPayslipService(repo)
    handler := NewPayslipHandler(service, log)

    routersGroup := router.Group("v1")
    {
        payslipGroup := routersGroup.Group("payslips")
        payslipGroup.Use(utils.JWTAuthMiddleware())
        payslipGroup.POST("/run", handler.RunPayroll)
		payslipGroup.POST("/generate", handler.GeneratePayslip)
    }
}