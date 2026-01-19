package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/example/go-api/utils"
)

// RequestID injects a unique request identifier into context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Request-Id")
		if traceID == "" {
			traceID = utils.NewID()
		}
		c.Set(utils.GinKeyTraceID, traceID)
		c.Writer.Header().Set("X-Request-Id", traceID)
		c.Request = c.Request.WithContext(utils.WithTraceID(c.Request.Context(), traceID))
		c.Next()
	}
}
