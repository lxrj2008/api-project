package controller

import "github.com/gin-gonic/gin"

// Health responds with OK when the service is running.
// @Summary Health Check
// @Description Returns service health status
// @Tags System
// @Produce json
// @Success 200 {object} APIResponse
// @Router /healthz [get]
func Health(c *gin.Context) {
	RespondSuccess(c, gin.H{"status": "ok"})
}
