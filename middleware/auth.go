package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-api/auth"
	"github.com/example/go-api/controller"
	"github.com/example/go-api/utils"
)

// Auth validates JWT bearer tokens.
func Auth(jwtManager *auth.JWTManager, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer"))
		if token == "" {
			respondUnauthorized(c)
			c.Abort()
			return
		}

		claims, err := jwtManager.Validate(token)
		if err != nil {
			logger.Warn("invalid token", zap.Error(err))
			respondUnauthorized(c)
			c.Abort()
			return
		}

		c.Set(utils.GinKeyUserID, claims.Subject)
		c.Request = c.Request.WithContext(utils.WithUserID(c.Request.Context(), claims.Subject))
		c.Next()
	}
}

func respondUnauthorized(c *gin.Context) {
	traceID := c.GetString(utils.GinKeyTraceID)
	c.JSON(http.StatusUnauthorized, controller.APIResponse{Code: utils.ErrUnauthorized.Code, Message: utils.ErrUnauthorized.Message, TraceID: traceID})
}
