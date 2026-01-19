package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-api/utils"
)

// APIResponse defines the unified response shape.
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Details interface{} `json:"details,omitempty"`
	TraceID string      `json:"traceId"`
}

// RespondSuccess writes a 200 response with data.
func RespondSuccess(c *gin.Context, data interface{}) {
	traceID := c.GetString(utils.GinKeyTraceID)
	c.JSON(http.StatusOK, APIResponse{Code: 0, Message: "OK", Data: data, TraceID: traceID})
}

// RespondMessage writes a simple OK response.
func RespondMessage(c *gin.Context, status int, message string) {
	traceID := c.GetString(utils.GinKeyTraceID)
	c.JSON(status, APIResponse{Code: 0, Message: message, TraceID: traceID})
}

// RespondError rewrites application errors with code + message.
func RespondError(c *gin.Context, logger *zap.Logger, err error) {
	traceID := c.GetString(utils.GinKeyTraceID)
	if appErr, ok := utils.IsAppError(err); ok {
		c.JSON(appErr.HTTPStatus, APIResponse{Code: appErr.Code, Message: appErr.Message, Details: appErr.Details, TraceID: traceID})
		return
	}

	logger.Error("unhandled error", zap.Error(err), zap.String("traceId", traceID))
	c.JSON(http.StatusInternalServerError, APIResponse{Code: utils.ErrInternal.Code, Message: utils.ErrInternal.Message, TraceID: traceID})
}
