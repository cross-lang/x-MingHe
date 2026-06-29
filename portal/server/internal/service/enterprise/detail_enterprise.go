package enterprise

import (
	"context"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/pkg/stringx"
	"portal/internal/types"

	"go.uber.org/zap"
)

// DetailEnterprise 查询企业详情
func (s *Service) DetailEnterprise(ctx context.Context, enterpriseId uint32) (*types.EnterpriseItem, error) {
	// 获取企业信息
	enterprise, err := s.EnterpriseRepo.GetById(ctx, enterpriseId)
	if err != nil {
		log.WithContext(ctx).Error("获取企业信息失败", zap.Error(err))
		return nil, errorx.FailedEnterpriseError
	}

	return &types.EnterpriseItem{
		ID:          enterprise.ID,
		Type:        enterprise.Type,
		FullName:    enterprise.FullName,
		ShortName:   enterprise.ShortName,
		Description: enterprise.Description,
		Icon:        enterprise.Icon,
		Images:      stringx.Split(enterprise.Image, ";"),
		Contact:     stringx.Split(enterprise.Contact, ";"),
		Website:     stringx.Split(enterprise.Website, ";"),
		Labels:      stringx.Split(enterprise.Label, ";"),
		Category:    make([]uint32, 0),
		LeaseTime:   enterprise.LeaseTime.Format("2006-01-02"),
		CreatedAt:   enterprise.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   enterprise.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
