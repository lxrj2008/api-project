package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"liangxiong/demo/dto"
	"liangxiong/demo/service"
	"liangxiong/demo/utils"
)

// UserController handles HTTP requests for users.
type UserController struct {
	service *service.UserService
	logger  *zap.Logger
}

// NewUserController constructs a controller.
func NewUserController(service *service.UserService, logger *zap.Logger) *UserController {
	return &UserController{service: service, logger: logger}
}

// List handles GET /users.
// @Summary List users
// @Description Paginated list of users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(20)
// @Success 200 {object} APIResponse{data=dto.UserListResponse}
// @Failure 401 {object} APIResponse
// @Router /api/v1/users [get]
func (ctl *UserController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	resp, err := ctl.service.ListUsers(c.Request.Context(), page, size)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondSuccess(c, resp)
}

// Get handles GET /users/:id.
// @Summary Get user
// @Description Retrieve a user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} APIResponse{data=dto.UserResponse}
// @Failure 401 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/v1/users/{id} [get]
func (ctl *UserController) Get(c *gin.Context) {
	id := c.Param("id")
	resp, err := ctl.service.GetUser(c.Request.Context(), id)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondSuccess(c, resp)
}

// Create handles POST /users.
// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UserCreateRequest true "User payload"
// @Success 201 {object} APIResponse{data=dto.UserResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Router /api/v1/users [post]
func (ctl *UserController) Create(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, ctl.logger, NewBindingError(err))
		return
	}
	resp, err := ctl.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	c.JSON(http.StatusCreated, APIResponse{Code: 0, Message: "Created", Data: resp, TraceID: c.GetString(utils.GinKeyTraceID)})
}

// Update handles PUT /users/:id.
// @Summary Update user
// @Description Update existing user fields
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UserUpdateRequest true "User payload"
// @Success 200 {object} APIResponse{data=dto.UserResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/v1/users/{id} [put]
func (ctl *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, ctl.logger, NewBindingError(err))
		return
	}
	resp, err := ctl.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondSuccess(c, resp)
}

// Delete handles DELETE /users/:id.
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/v1/users/{id} [delete]
func (ctl *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctl.service.DeleteUser(c.Request.Context(), id); err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondMessage(c, http.StatusOK, "Deleted")
}
