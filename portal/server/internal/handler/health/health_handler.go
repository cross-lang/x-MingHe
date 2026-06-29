package health

import (
	"net/http"
	"portal/internal/pkg/ginx"
	"portal/internal/service/health"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 健康检查处理器
type Handler struct {
	HealthService *health.Service
}

// NewHandler 创建健康检查处理器
func NewHandler(service *health.Service) *Handler {
	return &Handler{
		HealthService: service,
	}
}

// CheckHealth 健康检查接口
// @Summary 健康检查
// @Description 检查服务及其依赖组件的健康状态
// @Tags Health
// @Produce json
// @Param detailed query bool false "是否返回详细信息" default(false)
// @Success 200 {object} ginx.Response
// @Router /health [get]
func (h *Handler) CheckHealth(ctx *gin.Context) {
	// 获取详细模式参数
	detailedStr := ctx.DefaultQuery("detailed", "false")
	detailed, err := strconv.ParseBool(detailedStr)
	if err != nil {
		detailed = false
	}

	// 执行健康检查
	result, err := h.HealthService.CheckHealth(ctx.Request.Context(), detailed)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	// 根据状态返回不同的 HTTP 状态码
	statusCode := http.StatusOK
	if result.Status == string(health.HealthStatusUnhealthy) {
		statusCode = http.StatusServiceUnavailable
	} else if result.Status == string(health.HealthStatusDegraded) {
		statusCode = http.StatusOK // degraded 仍然返回 200，但状态为 degraded
	}

	// 返回响应
	ctx.JSON(statusCode, ginx.Response{
		Code: 0,
		Msg:  result.Status,
		Data: result,
	})
}

// Readiness 就绪检查接口
// @Summary 就绪检查
// @Description 检查服务是否准备好接收流量
// @Tags Health
// @Produce json
// @Success 200 {object} ginx.Response
// @Router /ready [get]
func (h *Handler) Readiness(ctx *gin.Context) {
	// 执行快速健康检查（不包含详细信息）
	result, err := h.HealthService.CheckHealth(ctx.Request.Context(), false)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	// 如果不是健康状态，返回 503
	statusCode := http.StatusOK
	if result.Status != string(health.HealthStatusHealthy) {
		statusCode = http.StatusServiceUnavailable
	}

	ctx.JSON(statusCode, ginx.Response{
		Code: 0,
		Msg:  "ready",
		Data: map[string]interface{}{
			"status": result.Status,
		},
	})
}

// Liveness 存活检查接口
// @Summary 存活检查
// @Description 检查服务是否正在运行
// @Tags Health
// @Produce json
// @Success 200 {object} ginx.Response
// @Router /live [get]
func (h *Handler) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ginx.Response{
		Code: 0,
		Msg:  "alive",
	})
}
