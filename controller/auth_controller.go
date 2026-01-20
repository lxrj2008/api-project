package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"liangxiong/demo/dto"
	"liangxiong/demo/service"
)

// AuthController exposes authentication endpoints.
type AuthController struct {
	service *service.AuthService
	logger  *zap.Logger
}

// NewAuthController builds the controller.
func NewAuthController(service *service.AuthService, logger *zap.Logger) *AuthController {
	return &AuthController{service: service, logger: logger}
}

// Login exchanges credentials for a JWT.
// @Summary Login
// @Description Authenticate a user and return JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} APIResponse{data=dto.LoginResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Router /api/v1/auth/login [post]
func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, ctl.logger, NewBindingError(err))
		return
	}

	resp, err := ctl.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}

	RespondSuccess(c, resp)
}
