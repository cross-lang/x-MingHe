package user

import (
	"context"
	"portal/internal/model"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// UserRegister 用户注册
func (s *Service) UserRegister(ctx context.Context, form *types.UserRegisterReq) (*types.UserLoginResp, error) {
	// 检查验证码是否正确
	code, err := s.Cache.GetRegisterCode(ctx, form.PhoneNumber)
	if err != nil {
		log.WithContext(ctx).Error("获取缓存中的验证码失败", zap.Error(err))
		return nil, errorx.VerificationCodeError
	}
	if code != form.Code {
		return nil, errorx.VerificationCodeError
	}
	// 清除验证码
	defer s.Cache.ClearRegisterCode(ctx, form.PhoneNumber)

	// 检查手机号是否已经注册
	exists, err := s.UserRepo.PhoneNumberExists(ctx, form.PhoneNumber)
	if err != nil {
		log.WithContext(ctx).Error("检查用户名是否已经注册失败", zap.Error(err))
		return nil, errorx.UserRegisterError
	}
	if exists {
		log.WithContext(ctx).Info("账号已经已经注册", zap.String("account", form.PhoneNumber))
		return nil, errorx.AccountRegisteredError
	}

	// 新增用户信息
	user := &model.XUser{
		Name:          form.PhoneNumber,
		BlockStatus:   1, // 默认未拉黑
		AccountStatus: 1, // 默认未禁用
		Account:       form.PhoneNumber,
		PhoneNumber:   form.PhoneNumber,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 保存用户注册信息
	if err = s.UserRepo.RegisterSave(ctx, user); err != nil {
		log.WithContext(ctx).Error("保存用户注册信息失败", zap.Error(err))
		return nil, errorx.UserRegisterError
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
	return &types.UserLoginResp{
		Token:     jwtToken,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}
