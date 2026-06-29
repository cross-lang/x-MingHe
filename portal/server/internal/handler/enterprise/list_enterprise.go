package enterprise

import (
	"portal/internal/pkg/ginx"
	"portal/internal/types"

	"github.com/gin-gonic/gin"
)

// ListEnterprise 获取指定园区的企业列表
// @Summary 获取指定园区的企业列表
// @Tags 企业
// @Param park_id path int true "园区ID"
// @Param type query int true "企业类型（1：企业，2：学院）"
// @Param category query int false "类目"
// @Param size query int true "每页数据量"
// @Param page query int true "页码"
// @Param keyword query string false "关键字"
// @Produce json
// @Success 200 {object} ginx.Response{data=types.EnterpriseList} "获取指定园区的企业列表成功"
// @Security ApiKeyAuth
// @Router /v1/enterprises/list [get]
func (h *Handler) ListEnterprise(ctx *gin.Context) {
	var form types.ListEnterpriseReq
	if err := ginx.ShouldBindUrl(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	enterpriseList, err := h.EnterpriseService.ListEnterprise(ctx.Request.Context(), &form)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, enterpriseList)
}
