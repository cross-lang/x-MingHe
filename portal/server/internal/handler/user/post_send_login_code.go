package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

// PostSendLoginCode 发送登录短信验证码
// @Summary 发送登录短信验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body PostSendSMSCodeReq true "参数"
// @Success 200 {object} ginx.Response{data=string} "短信验证码发送成功"
// @Router /v1/users/login/code/send [post]
func (h *Handler) PostSendLoginCode(ctx *gin.Context) {
	var form PostSendSMSCodeReq
	if err := ginx.ShouldBind(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	// 发送验证码
	err := h.UserService.SendLoginCode(ctx, form.PhoneNumber)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, "短信验证码发送成功")
}
