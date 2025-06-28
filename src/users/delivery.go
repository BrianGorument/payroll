package users

import (
	"net/http"
	"payroll/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserHandler struct
type UserHandler struct {
	service IUserService
	logger  *logrus.Logger
}

// NewUserHandler (Dependency Injection)
func NewUserHandler(service IUserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{service, logger}
}

// Endpoint: Get All Users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		h.logger.Error("Failed to get users:", err)
		response.ErrorHandler(c, h.logger, nil, err)
		return
	}

	if len(users) == 0 {
		resp := response.ErrorStruct{
			HTTPCode:           http.StatusNotFound,
			Description:        response.DescriptionFailed,
			Message:            response.DataNotFound,
			MessageDescription: "Users data is empty and not found",
		}
		response.SendErrorResponse(c, http.StatusNotFound, resp)
		return
	}
	succesresp := response.Response{
		Description:        response.DescriptionSuccess,
		Message:            response.DataSuccess,
		MessageDescription: "Data retrieved successfully",
		Data:               users,
	}
	response.SendResponseSuccess(c, http.StatusOK, succesresp)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body request"})
		return
	}

	user, err := h.service.LoginUser(req)
	if err != nil {
		h.logger.Error("Failed to login:", err)
		resp := response.ErrorStruct{
			Description:        response.DescriptionFailed,
			Message:            err.Error(),
			MessageDescription: "Failed to login",
			Data:               err,
		}
		response.SendErrorResponse(c, http.StatusBadRequest, resp)
		return
	}

	succesresp := response.Response{
		Description:        response.DescriptionSuccess,
		Message:            response.DataSuccess,
		MessageDescription: "Successfully logged in",
		Data:               user,
	}
	response.SendResponseSuccess(c, http.StatusOK, succesresp)
}
