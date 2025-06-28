package overtimes

import (
    "net/http"
    "payroll/shared/response"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type OvertimeHandler struct {
    service IOvertimeService
    logger  *logrus.Logger
}

func NewOvertimeHandler(service IOvertimeService, logger *logrus.Logger) *OvertimeHandler {
    return &OvertimeHandler{service, logger}
}

func (h *OvertimeHandler) CreateOvertime(c *gin.Context) {
    var req CreateOvertimeRequest
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

    // Validasi role (hanya employee)
    role, exists := c.Get("role")
    if !exists || role != "employee" {
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusForbidden,
            Description:        response.DescriptionFailed,
            Message:            "Forbidden",
            MessageDescription: "Only employees can submit overtime",
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

    overtime, err := h.service.CreateOvertime(req, userIDStr)
    if err != nil {
        h.logger.Error("Failed to create overtime:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusBadRequest,
            Description:        response.DescriptionFailed,
            Message:            err.Error(),
            MessageDescription: "Failed to create overtime",
            Data:               err,
        }
        response.SendErrorResponse(c, http.StatusBadRequest, resp)
        return
    }

    // Set untuk AuditLogMiddleware
    c.Set("record_id", overtime.ID)
    c.Set("response_data", overtime)

    succesresp := response.Response{
        Description:        response.DescriptionSuccess,
        Message:            response.DataSuccess,
        MessageDescription: response.SuccessInsert,
        Data:               overtime,
    }
    response.SendResponseSuccess(c, http.StatusCreated, succesresp)
}