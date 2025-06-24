package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Response model
type Response struct {
	ResponseCode       string      `json:"rc"`
	Description        string      `json:"desc"`
	Message            string      `json:"msg"`
	MessageDescription string      `json:"msg_desc"`
	Data               interface{} `json:"data,omitempty"`
}

// SendSuccessResponse mengirimkan response sukses
func SendSuccessResponse(ctx *gin.Context, data interface{}) {
	resp := Response{
		ResponseCode:       RCSuccess,
		Description:        DescriptionSuccess,
		Message:            DataSuccess,
		MessageDescription: DataSuccessDesc,
		Data:               data,
	}
	ctx.JSON(http.StatusOK, resp)
}

// SendErrorResponse mengirimkan response error dengan kode HTTP yang sesuai
func SendErrorResponse(ctx *gin.Context, httpCode int, errStruct ErrorStruct) {
	ctx.JSON(httpCode, errStruct)
}

func SendResponseSuccess(ctx *gin.Context, httpCode int, SuccessStruct Response) {
	ctx.JSON(httpCode, SuccessStruct)
}

// ErrorHandler menangani error dalam controller
func ErrorHandler(ctx *gin.Context, logger *logrus.Logger, req interface{}, err error) {
	logger.Error("Error processing request:", err)

	resp := ErrorStruct{
		HTTPCode:           http.StatusInternalServerError,
		Code:               RCServerError,
		Description:        DescriptionFailed,
		Message:            ServerError,
		MessageDescription: ServerErrorDesc,
	}

	if customErr, ok := err.(ErrorStruct); ok {
		resp = customErr
	}

	SendErrorResponse(ctx, resp.HTTPCode, resp)
}
