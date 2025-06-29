package payslips

import (
	"net/http"
	"os"
	"path/filepath"
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
func (h *PayslipHandler) GeneratePayslip(c *gin.Context) {
    var req GeneratePayslipRequest
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

    // Validasi role (admin atau employee)
    role, exists := c.Get("role")
    if !exists || (role != "admin" && role != "employee") {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusForbidden,
            Description:        response.DescriptionFailed,
            Message:            "Forbidden",
            MessageDescription: "Only admins or employees can generate payslip",
        }
        response.SendErrorResponse(c, http.StatusForbidden, resp)
        return
    }

    // Ambil userID dari JWT
    userID, exists := c.Get("userid")
    if !exists || userID == "" {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusUnauthorized,
            Description:        response.DescriptionFailed,
            Message:            "Unauthorized",
            MessageDescription: "User ID not found",
        }
        response.SendErrorResponse(c, http.StatusUnauthorized, resp)
        return
    }

    userIDStr, ok := userID.(string)
    if !ok {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusInternalServerError,
            Description:        response.DescriptionFailed,
            Message:            "Internal Server Error",
            MessageDescription: "Invalid user ID format",
        }
        response.SendErrorResponse(c, http.StatusInternalServerError, resp)
        return
    }

    // Generate payslip PDF
    pdfPath, err := h.service.GeneratePayslip(req, userIDStr, role.(string))
    if err != nil {
        h.logger.Error("Failed to generate payslip:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusBadRequest,
            Description:        response.DescriptionFailed,
            Message:            err.Error(),
            MessageDescription: "Failed to generate payslip",
            Data:               err,
        }
        response.SendErrorResponse(c, http.StatusBadRequest, resp)
        return
    }

    // Baca file PDF
    pdfData, err := os.ReadFile(pdfPath)
    if err != nil {
        h.logger.Error("Failed to read PDF file:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusInternalServerError,
            Description:        response.DescriptionFailed,
            Message:            "Failed to read payslip file",
            MessageDescription: err.Error(),
        }
        response.SendErrorResponse(c, http.StatusInternalServerError, resp)
        return
    }

    // Set header untuk streaming PDF
    c.Header("Content-Type", "application/pdf")
    c.Header("Content-Disposition", "attachment; filename="+filepath.Base(pdfPath))
    c.Data(http.StatusOK, "application/pdf", pdfData)
}