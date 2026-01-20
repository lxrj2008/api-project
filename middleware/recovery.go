package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"liangxiong/demo/controller"
	"liangxiong/demo/utils"
)

// Recovery handles panics and emits JSON error.
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", zap.Any("panic", r))
				c.AbortWithStatusJSON(http.StatusInternalServerError, controller.APIResponse{Code: utils.ErrInternal.Code, Message: utils.ErrInternal.Message, TraceID: c.GetString(utils.GinKeyTraceID)})
			}
		}()
		c.Next()
	}
}
