package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

// PostDeactivateUser 用户注销
// @Summary 用户注销接口
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} ginx.Response{data=types.DetailDeactivateUserCheckResp} "用户注销成功"
// @Security ApiKeyAuth
// @Router /v1/users/deactivate [post]
func (h *Handler) PostDeactivateUser(ctx *gin.Context) {
	resp, err := h.UserService.DeactivateUser(ctx.Request.Context())
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	ginx.SuccessResponse(ctx, resp)
}
