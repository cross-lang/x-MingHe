package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

// DetailCheckDeactivate 查询用户注销准入评估结果
// @Summary 查询用户注销准入评估结果接口
// @Tags 用户
// @Produce json
// @Success 200 {object} ginx.Response{data=types.DetailDeactivateUserCheckResp} "查询用户注销准入评估结果成功"
// @Security ApiKeyAuth
// @Router /v1/users/check/deactivate/detail [get]
func (h *Handler) DetailCheckDeactivate(ctx *gin.Context) {
	resp, err := h.UserService.DetailDeactivateUserCheck(ctx.Request.Context())
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	ginx.SuccessResponse(ctx, resp)
}
