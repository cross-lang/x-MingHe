package enterprise

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

type DetailEnterpriseReq struct {
	EnterpriseId uint32 `uri:"enterprise_id" binding:"required"`
}

// DetailEnterprise 查询企业详情
// @Summary 查询企业详情
// @Tags 企业
// @Param enterprise_id path int true "企业ID"
// @Produce json
// @Success 200 {object} ginx.Response{data=types.EnterpriseItem} "查询企业详情成功"
// @Security ApiKeyAuth
// @Router /v1/enterprises/{enterprise_id}/detail [get]
func (h *Handler) DetailEnterprise(ctx *gin.Context) {
	var form DetailEnterpriseReq
	if err := ginx.ShouldBindUrl(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	info, err := h.EnterpriseService.DetailEnterprise(ctx.Request.Context(), form.EnterpriseId)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, info)
}
