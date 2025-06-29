package payslips

import (
    "net/http"
    "payroll/shared/response"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type PayslipHandler struct {
    service IPayslipService
    logger  *logrus.Logger
}

func NewPayslipHandler(service IPayslipService, logger *logrus.Logger) *PayslipHandler {
    return &PayslipHandler{service, logger}
}

func (h *PayslipHandler) RunPayroll(c *gin.Context) {
    var req RunPayrollRequest
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

    // Validasi role (hanya admin)
    role, exists := c.Get("role")
    if !exists || role != "admin" {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusForbidden,
            Description:        response.DescriptionFailed,
            Message:            "Forbidden",
            MessageDescription: "Only admins can run payroll",
        }
        response.SendErrorResponse(c, http.StatusForbidden, resp)
        return
    }

    // Ambil adminID dari JWT
    adminID, exists := c.Get("userid")
    if !exists || adminID == "" {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusUnauthorized,
            Description:        response.DescriptionFailed,
            Message:            "Unauthorized",
            MessageDescription: "Admin ID not found",
        }
        response.SendErrorResponse(c, http.StatusUnauthorized, resp)
        return
    }

    adminIDStr, ok := adminID.(string)
    if !ok {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusInternalServerError,
            Description:        response.DescriptionFailed,
            Message:            "Internal Server Error",
            MessageDescription: "Invalid admin ID format",
        }
        response.SendErrorResponse(c, http.StatusInternalServerError, resp)
        return
    }

    payslips, err := h.service.RunPayroll(req, adminIDStr)
    if err != nil {
        h.logger.Error("Failed to run payroll:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusBadRequest,
            Description:        response.DescriptionFailed,
            Message:            err.Error(),
            MessageDescription: "Failed to run payroll",
            Data:               err,
        }
        response.SendErrorResponse(c, http.StatusBadRequest, resp)
        return
    }

    // Set untuk AuditLogMiddleware
    c.Set("record_id", "multiple_payslips")
    c.Set("response_data", payslips)

    succesresp := response.Response{
        Description:        response.DescriptionSuccess,
        Message:            response.DataSuccess,
        MessageDescription: "Payroll processed successfully",
        Data:               payslips,
    }
    response.SendResponseSuccess(c, http.StatusCreated, succesresp)
}