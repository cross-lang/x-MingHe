package user

import (
	"context"
	"portal/internal/constant"
	"portal/internal/model"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// UserLogin 用户登录
func (s *Service) UserLogin(ctx context.Context, form *types.UserLoginReq) (*types.UserLoginResp, error) {
	// 检查验证码是否正确
	code, err := s.Cache.GetLoginCode(ctx, form.PhoneNumber)
	if err != nil {
		log.WithContext(ctx).Error("获取缓存中的验证码失败", zap.Error(err))
		return nil, errorx.VerificationCodeError
	}
	if code != form.Code {
		return nil, errorx.VerificationCodeError
	}
	// 清理验证码
	defer s.Cache.ClearLoginCode(ctx, form.PhoneNumber)

	// 获取账号对应的用户信息
	user, err := s.UserRepo.GetByPhoneNumber(ctx, form.PhoneNumber)
	if err != nil {
		log.WithContext(ctx).Error("获取账号对应的用户信息", zap.Error(err))
		return nil, errorx.UserLoginError
	}

	if user.AccountStatus == constant.UserAccountStatusDisable {
		return nil, errorx.AccountDisableError
	}

	// 生成jwt返回
	expiresAt := time.Now().Add(s.Config.LoginJwt.Expires)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.UserLoginJwtClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.Config.LoginJwt.Issuer,
		},
	})
	secret := []byte(s.Config.LoginJwt.Key)
	jwtToken, err := token.SignedString(secret)
	if err != nil {
		log.WithContext(ctx).Error("生成jwt失败", zap.Error(err))
		return nil, errorx.UserLoginError
	}

	// 设置用户登录状态
	err = s.Cache.SetLoginStatus(ctx, user.ID, s.Config.LoginJwt.Expires, jwtToken)
	if err != nil {
		log.WithContext(ctx).Error("设置用户登录状态失败", zap.Error(err))
		return nil, errorx.UserLoginError
	}

	// 检查用户是否在注销中，如果状态是注销中就恢复为正常状态
	s.DeactivateRecovery(ctx, user)

	return &types.UserLoginResp{
		Token:     jwtToken,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// DeactivateRecovery 检查用户是否在注销中，如果状态是注销中就恢复为正常状态
func (s *Service) DeactivateRecovery(ctx context.Context, user *model.XUser) {
	// 只有当用户状态为注销中时才需要恢复
	if user.DeactivateStatus != constant.DeactivateUserStatusPending {
		return
	}

	// 将注销状态恢复为未注销，注销时间清空
	err := s.UserRepo.MarkDeactivateStatus(ctx, user.ID, constant.DeactivateUserStatusNot, nil)
	if err != nil {
		log.WithContext(ctx).Error("登录时恢复用户注销状态失败", zap.Error(err))
		return
	}
}
