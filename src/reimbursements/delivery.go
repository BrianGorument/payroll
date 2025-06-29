package reimbursements

import (
    "net/http"
    "payroll/shared/response"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type ReimbursementHandler struct {
    service IReimbursementService
    logger  *logrus.Logger
}

func NewReimbursementHandler(service IReimbursementService, logger *logrus.Logger) *ReimbursementHandler {
    return &ReimbursementHandler{service, logger}
}

func (h *ReimbursementHandler) CreateReimbursement(c *gin.Context) {
    var req CreateReimbursementRequest
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
            MessageDescription: "Only employees can submit reimbursement",
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

    reimbursement, err := h.service.CreateReimbursement(req, userIDStr)
    if err != nil {
        h.logger.Error("Failed to create reimbursement:", err)
        resp := response.ErrorStruct{
            HTTPCode:           http.StatusBadRequest,
            Description:        response.DescriptionFailed,
            Message:            err.Error(),
            MessageDescription: "Failed to create reimbursement",
            Data:               err,
        }
        response.SendErrorResponse(c, http.StatusBadRequest, resp)
        return
    }

    // Set untuk AuditLogMiddleware
    c.Set("record_id", reimbursement.ID)
    c.Set("response_data", reimbursement)

    succesresp := response.Response{
        Description:        response.DescriptionSuccess,
        Message:            response.DataSuccess,
        MessageDescription: response.SuccessInsert,
        Data:               reimbursement,
    }
    response.SendResponseSuccess(c, http.StatusCreated, succesresp)
}