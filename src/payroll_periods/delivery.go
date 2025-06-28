package payroll_periods

import (
    "net/http"
    "payroll/shared/response"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
	"fmt"
)

type PayrollPeriodHandler struct {
    service IPayrollPeriodService
    logger  *logrus.Logger
}

func NewPayrollPeriodHandler(service IPayrollPeriodService, logger *logrus.Logger) *PayrollPeriodHandler {
    return &PayrollPeriodHandler{service, logger}
}

func (h *PayrollPeriodHandler) CreatePayrollPeriod(c *gin.Context) {
    var req CreatePayrollPeriodRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("Invalid request:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusBadRequest,
            Description:        response.DescriptionFailed,
            Message:            "Invalid Body request",
            MessageDescription: err.Error(),
        }
        response.SendErrorResponse(c, http.StatusBadRequest, resp)
        return
    }

    role, exists := c.Get("role")
    if !exists || role != "admin" {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusForbidden,
            Description:        response.DescriptionFailed,
            Message:            "Forbidden",
            MessageDescription: "Only admins can create payroll periods",
        }
        response.SendErrorResponse(c, http.StatusForbidden, resp)
        return
    }

    userID, exists := c.Get("userId")
    if !exists {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusUnauthorized,
            Description:        response.DescriptionFailed,
            Message:            "Unauthorized",
            MessageDescription: "User ID not found",
        }
        response.SendErrorResponse(c, http.StatusUnauthorized, resp)
        return
    }
	
	fmt.Println(userID)

    period, err := h.service.CreatePayrollPeriod(req, userID.(string))
    if err != nil {
        h.logger.Error("Failed to create payroll period:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusBadRequest,
            Description:        response.DescriptionFailed,
            Message:            err.Error(),
            MessageDescription: "Failed to create payroll period",
            Data:               err,
        }
        response.SendErrorResponse(c, http.StatusBadRequest, resp)
        return
    }

    c.Set("record_id", period.ID)
    c.Set("response_data", period)

    succesresp := response.Response{
        Description:        response.DescriptionSuccess,
        Message:            response.DataSuccess,
        MessageDescription: response.SuccessInsert,
        Data:               period,
    }
    response.SendResponseSuccess(c, http.StatusCreated, succesresp)
}