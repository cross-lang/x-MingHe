package enterprise

import (
	"context"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/pkg/stringx"
	"portal/internal/repository"
	"portal/internal/types"

	"go.uber.org/zap"
)

// ListEnterprise 获取指定园区下的企业列表
func (s *Service) ListEnterprise(ctx context.Context, form *types.ListEnterpriseReq) (*types.EnterpriseList, error) {
	// 获取园区的企业列表
	enterprises, count, err := s.EnterpriseRepo.SearchEnterpriseList(ctx, &repository.SearchEnterpriseListParams{
		ParkId:   form.ParkId,
		Type:     form.Type,
		Keyword:  form.Keyword,
		Category: form.Category,
		Size:     form.Size,
		Page:     form.Page,
	})
	if err != nil {
		log.WithContext(ctx).Error("获取园区的企业列表失败", zap.Error(err))
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

	return &types.EnterpriseList{
		Count: count,
		Items: items,
	}, nil
}
