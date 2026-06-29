package user

import (
	"portal/internal/pkg/ginx"
	"portal/internal/types"

	"github.com/gin-gonic/gin"
)

// PostVerifyUser 用户实名认证
// @Summary 用户实名认证
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body types.UserVerificationReq true "用户注册请求参数"
// @Success 200 {object} ginx.Response{data=string} "用户实名认证成功"
// @Security ApiKeyAuth
// @Router /v1/users/verification [post]
func (h *Handler) PostVerifyUser(ctx *gin.Context) {
	var form types.UserVerificationReq
	if err := ginx.ShouldBind(ctx, &form); err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	user, err := ginx.DetailUserFromCtx(ctx)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}

	err = h.UserService.UserVerification(ctx, user, &form)
	if err != nil {
		ginx.FailureResponse(ctx, err)
		return
	}
	ginx.SuccessResponse(ctx, "success")
}
