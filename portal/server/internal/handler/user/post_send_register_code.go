package user

import (
	"portal/internal/pkg/ginx"

	"github.com/gin-gonic/gin"
)

type PostSendSMSCodeReq struct {
	PhoneNumber string `json:"phone_number" binding:"required"` // 账号
}

// PostSendRegisterCode 发送注册短信验证码
// @Summary 发送注册短信验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body PostSendSMSCodeReq true "参数"
// @Success 200 {object} ginx.Response{data=string} "短信验证码发送成功"
// @Router /v1/users/register/code/send [post]
func (h *Handler) PostSendRegisterCode(ctx *gin.Context) {
	var form PostSendSMSCodeReq
	if err := ginx.ShouldBind(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	// 发送验证码
	err := h.UserService.SendRegisterCode(ctx, form.PhoneNumber)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, "短信验证码发送成功")
}
