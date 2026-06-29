package enterprise

import (
	"context"
	"portal/internal/model"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/types"

	"go.uber.org/zap"
)

// JoinEnterprise 加入企业
func (s *Service) JoinEnterprise(ctx context.Context, user *model.XUser, form *types.PostJoinEnterpriseReq) error {
	if user.VerificationStatus == 0 {
		// TODO: 需要实名认证才能加入企业
		// return errorx.UserNoVerifyError
	}

	// 获取企业信息
	enterprise, err := s.EnterpriseRepo.GetById(ctx, form.EnterpriseId)
	if err != nil {
		log.WithContext(ctx).Error("获取企业信息失败", zap.Error(err))
		return errorx.FailedEnterpriseError
	}

	// 检查用户是否已经入职这个企业
	has, err := s.UserEnterpriseRepo.UserHasEnterprise(ctx, user.ID, enterprise.ID)
	if err != nil {
		log.WithContext(ctx).Error("检查用户是否已经入职这个企业失败", zap.Error(err))
		return errorx.JoinEnterpriseError
	}
	if has {
		return errorx.UserHasEnterpriseError
	}

	// 用户企业关系
	userEnterprise := model.XUserEnterprise{
		UID:  user.ID,
		EID:  enterprise.ID,
		PID:  enterprise.PID,
		Type: enterprise.Type,
	}

	// 保存用户企业管理关系
	if err := s.UserEnterpriseRepo.SaveUserEnterprise(ctx, &userEnterprise); err != nil {
		log.WithContext(ctx).Error("保存用户入职企业关系失败", zap.Error(err))
		return errorx.JoinEnterpriseError
	}

	return nil
}
