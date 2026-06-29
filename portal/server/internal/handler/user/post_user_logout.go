package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

// PostLogoutUser 用户登出
// @Summary 用户登出接口
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} ginx.Response{data=string} "登出成功"
// @Security ApiKeyAuth
// @Router /v1/users/logout [post]
func (h *Handler) PostLogoutUser(ctx *gin.Context) {
	userInfo, err := ginx.DetailUserFromCtx(ctx)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	err = h.UserService.UserLogout(ctx, userInfo.ID)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, "success")
}
