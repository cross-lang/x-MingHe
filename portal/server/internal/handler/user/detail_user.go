package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

// DetailUser 查询用户详情
// @Summary 查询用户详情接口
// @Tags 用户
// @Produce json
// @Success 200 {object} ginx.Response{data=types.UserItem} "查询用户详情成功"
// @Security ApiKeyAuth
// @Router /v1/users/detail [get]
func (h *Handler) DetailUser(ctx *gin.Context) {
	detailUser, err := h.UserService.DetailUser(ctx.Request.Context())
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, detailUser)
}
