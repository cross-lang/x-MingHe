package user

import (
	"portal/internal/pkg/ginx"
	"portal/internal/types"

	"github.com/gin-gonic/gin"
)

// PostLoginUser 用户登录
// @Summary 用户登录接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body types.UserLoginReq true "用户登录请求参数"
// @Success 200 {object} ginx.Response{data=types.UserLoginResp} "登录成功"
// @Router /v1/users/login [post]
func (h *Handler) PostLoginUser(ctx *gin.Context) {
	var form types.UserLoginReq
	if err := ginx.ShouldBind(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	jwtToken, err := h.UserService.UserLogin(ctx, &form)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, jwtToken)
}
