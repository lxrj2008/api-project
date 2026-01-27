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

// ExchangeController handles exchange endpoints.
type ExchangeController struct {
	service *service.ExchangeService
	logger  *zap.Logger
}

// NewExchangeController constructs the controller.
func NewExchangeController(service *service.ExchangeService, logger *zap.Logger) *ExchangeController {
	return &ExchangeController{service: service, logger: logger}
}

// List handles GET /exchanges.
// @Summary List exchanges
// @Description Paginated list of exchanges
// @Tags Exchanges
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(20)
// @Success 200 {object} APIResponse{data=dto.ExchangeListResponse}
// @Failure 401 {object} APIResponse
// @Router /api/v1/exchanges [get]
func (ctl *ExchangeController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	resp, err := ctl.service.ListExchanges(c.Request.Context(), page, size)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondSuccess(c, resp)
}

// Get handles GET /exchanges/:code.
// @Summary Get exchange
// @Description Retrieve an exchange by MQM code
// @Tags Exchanges
// @Security BearerAuth
// @Produce json
// @Param code path string true "MQM Exchange Code"
// @Success 200 {object} APIResponse{data=dto.ExchangeResponse}
// @Failure 401 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/v1/exchanges/{code} [get]
func (ctl *ExchangeController) Get(c *gin.Context) {
	code := c.Param("code")
	resp, err := ctl.service.GetExchange(c.Request.Context(), code)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondSuccess(c, resp)
}

// Create handles POST /exchanges.
// @Summary Create exchange
// @Description Create a new exchange
// @Tags Exchanges
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ExchangeCreateRequest true "Exchange payload"
// @Success 201 {object} APIResponse{data=dto.ExchangeResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Router /api/v1/exchanges [post]
func (ctl *ExchangeController) Create(c *gin.Context) {
	var req dto.ExchangeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, ctl.logger, NewBindingError(err))
		return
	}
	resp, err := ctl.service.CreateExchange(c.Request.Context(), req)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	c.JSON(http.StatusCreated, APIResponse{Code: 0, Message: "Created", Data: resp, TraceID: c.GetString(utils.GinKeyTraceID)})
}

// Update handles PUT /exchanges/:code.
// @Summary Update exchange
// @Description Update an exchange record
// @Tags Exchanges
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code path string true "MQM Exchange Code"
// @Param request body dto.ExchangeUpdateRequest true "Exchange payload"
// @Success 200 {object} APIResponse{data=dto.ExchangeResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/v1/exchanges/{code} [put]
func (ctl *ExchangeController) Update(c *gin.Context) {
	code := c.Param("code")
	var req dto.ExchangeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, ctl.logger, NewBindingError(err))
		return
	}
	resp, err := ctl.service.UpdateExchange(c.Request.Context(), code, req)
	if err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondSuccess(c, resp)
}

// Delete handles DELETE /exchanges/:code.
// @Summary Delete exchange
// @Description Delete an exchange by MQM code
// @Tags Exchanges
// @Security BearerAuth
// @Produce json
// @Param code path string true "MQM Exchange Code"
// @Success 200 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /api/v1/exchanges/{code} [delete]
func (ctl *ExchangeController) Delete(c *gin.Context) {
	code := c.Param("code")
	if err := ctl.service.DeleteExchange(c.Request.Context(), code); err != nil {
		RespondError(c, ctl.logger, err)
		return
	}
	RespondMessage(c, http.StatusOK, "Deleted")
}
