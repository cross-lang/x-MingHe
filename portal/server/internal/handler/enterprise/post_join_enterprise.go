package enterprise

import (
	"portal/internal/pkg/ginx"
	"portal/internal/types"

	"github.com/gin-gonic/gin"
)

// PostJoinEnterprise 用户加入企业
// @Summary 用户加入企业
// @Tags 企业
// @Accept json
// @Param data body types.PostJoinEnterpriseReq true "用户加入企业参数"
// @Produce json
// @Success 200 {object} ginx.Response{data=string} "用户加入企业成功"
// @Security ApiKeyAuth
// @Router /v1/enterprises/join [post]
func (h *Handler) PostJoinEnterprise(ctx *gin.Context) {
	var form types.PostJoinEnterpriseReq
	if err := ginx.ShouldBind(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	user, err := ginx.DetailUserFromCtx(ctx)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	err = h.EnterpriseService.JoinEnterprise(ctx.Request.Context(), user, &form)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, "success")
}
