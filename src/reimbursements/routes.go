package reimbursements

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "payroll/shared/utils"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, log *logrus.Logger) {
    repo := NewReimbursementRepository(db)
    service := NewReimbursementService(repo)
    handler := NewReimbursementHandler(service, log)

    routersGroup := router.Group("v1")
    {
        reimbursementGroup := routersGroup.Group("reimbursements")
        reimbursementGroup.Use(utils.JWTAuthMiddleware())
        reimbursementGroup.POST("/create", handler.CreateReimbursement)
    }
}