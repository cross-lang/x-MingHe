package user

import (
	"context"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/pkg/stringx"
	"portal/internal/types"

	"go.uber.org/zap"
)

// ListUserEnterprise 查询用户入驻的企业列表
func (s *Service) ListUserEnterprise(ctx context.Context, userId uint32) ([]*types.EnterpriseItem, error) {
	// 查询用户入驻的企业列表
	enterprises, err := s.UserEnterpriseRepo.GetUserEnterprises(ctx, userId)
	if err != nil {
		log.WithContext(ctx).Error("获取用户入驻的企业列表失败", zap.Error(err))
		return nil, errorx.FailedEnterpriseListError
	}

	// 数据处理
	items := make([]*types.EnterpriseItem, 0, len(enterprises))
	for _, item := range enterprises {
		items = append(items, &types.EnterpriseItem{
			ID:        item.ID,
			Type:      item.Type,
			FullName:  item.FullName,
			ShortName: item.ShortName,
			Icon:      item.Icon,
			Labels:    stringx.Split(item.Label, ";"),
			Category:  make([]uint32, 0),
			Top:       item.Top,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return items, nil
}
