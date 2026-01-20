package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"liangxiong/demo/controller"
	"liangxiong/demo/utils"
)

// RateLimit limits per-process requests per second.
func RateLimit(rps int) gin.HandlerFunc {
	if rps <= 0 {
		rps = 100
	}
	limiter := rate.NewLimiter(rate.Limit(rps), rps)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			traceID := c.GetString(utils.GinKeyTraceID)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, controller.APIResponse{Code: 429001, Message: "Too Many Requests", TraceID: traceID})
			return
		}
		c.Next()
	}
}
