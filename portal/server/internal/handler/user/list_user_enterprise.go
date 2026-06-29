package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

// ListUserEnterprise 查询用户入驻的企业列表
// @Summary 查询用户入驻的企业列表
// @Tags 用户
// @Produce json
// @Success 200 {object} ginx.Response{data=[]types.EnterpriseItem} "获取用户入驻的企业列表成功"
// @Security ApiKeyAuth
// @Router /v1/users/enterprises/list [get]
func (h *Handler) ListUserEnterprise(ctx *gin.Context) {
	user, err := ginx.DetailUserFromCtx(ctx)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	enterprises, err := h.UserService.ListUserEnterprise(ctx, user.ID)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, enterprises)
}
