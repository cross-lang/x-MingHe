package user

import (
	"portal/internal/pkg/ginx"
	"portal/internal/types"

	"github.com/gin-gonic/gin"
)

// PostRegisterUser 用户注册
// @Summary 用户注册接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body types.UserRegisterReq true "用户注册请求参数"
// @Success 200 {object} ginx.Response{data=types.UserLoginResp} "用户注册成功"
// @Router /v1/users/register [post]
func (h *Handler) PostRegisterUser(ctx *gin.Context) {
	var form types.UserRegisterReq
	if err := ginx.ShouldBind(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	jwtToken, err := h.UserService.UserRegister(ctx, &form)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, jwtToken)
}
