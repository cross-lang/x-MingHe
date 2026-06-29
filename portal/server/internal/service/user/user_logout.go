package user

import (
	"context"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"

	"go.uber.org/zap"
)

// UserLogout 用户登出
func (s *Service) UserLogout(ctx context.Context, userId uint32) error {

	// 设置用户登录状态
	err := s.Cache.ClearLoginStatus(ctx, userId)
	if err != nil {
		log.WithContext(ctx).Error("设置用户登录状态失败", zap.Error(err))
		return errorx.UserLogoutError
	}
	return nil
}
