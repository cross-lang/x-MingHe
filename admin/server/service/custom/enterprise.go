package custom

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type EnterpriseService struct{}

const (
	EnterpriseType = 1 // 企业
	CampusType     = 2 // 学院
)

// CreateEnterprise 创建企业（仅操作 x_enterprise 表）
func (s *EnterpriseService) CreateEnterprise(ctx context.Context, userId uint, params request.CreateEnterpriseParams) (any, error) {
	// 示例项目：使用第一个入驻时间
	leaseTime, _ := time.Parse("2006-01-02", params.Lease[0].LeaseTime)
	enterprise := &custom.XEnterprise{
		PID:         params.PID,
		Type:        EnterpriseType,
		FullName:    params.FullName,
		ShortName:   params.ShortName,
		Description: params.Description,
		Icon:        params.Icon,
		Image:       strings.Join(params.Images, ";"),
		Contact:     strings.Join(params.Contact, ";"),
		Website:     strings.Join(params.Websites, ";"),
		Label:       strings.Join(params.Labels, ";"),
		LeaseTime:   leaseTime,
		CreatedBy:   uint32(userId),
		UpdatedBy:   uint32(userId),
	}

	err := global.GVA_DB.WithContext(ctx).Create(enterprise).Error
	if err != nil {
		return nil, fmt.Errorf("保存企业信息失败，%v", err)
	}
	return enterprise.ID, nil
}

// ListEnterprise 获取企业列表（仅操作 x_enterprise、x_user_enterprise、x_user 表）
func (s *EnterpriseService) ListEnterprise(ctx context.Context, params request.ListEnterpriseParams) (any, int64, error) {
	query := global.GVA_DB.WithContext(ctx).Model(&custom.XEnterprise{}).Where("type = ?", EnterpriseType)

	// 园区过滤
	if params.PID != 0 {
		query = query.Where("p_id = ?", params.PID)
	}

	// 名称搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("full_name LIKE ? OR short_name LIKE ?", keyword, keyword)
	}

	// 入驻时间范围
	if params.LeaseStartTime != "" {
		query = query.Where("lease_time >= ?", params.LeaseStartTime)
	}
	if params.LeaseEndTime != "" {
		query = query.Where("lease_time <= ?", params.LeaseEndTime)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取企业总数失败，%v", err)
	}

	// 分页
	limit := params.PageSize
	offset := params.PageSize * (params.Page - 1)

	// 排序
	orderStr := "top DESC, weight DESC, updated_at DESC"

	// 查询企业数据
	var enterprises []custom.XEnterprise
	if err := query.Order(orderStr).Limit(limit).Offset(offset).Find(&enterprises).Error; err != nil {
		return nil, 0, fmt.Errorf("获取企业列表失败，%v", err)
	}

	// 获取企业ID列表
	eIds := make([]uint32, 0, len(enterprises))
	for _, e := range enterprises {
		eIds = append(eIds, e.ID)
	}

	// 查询企业关联的用户数（x_user_enterprise 表）
	userEnterpriseStats := make([]struct {
		EID   uint32
		Count int64
	}, 0)
	global.GVA_DB.WithContext(ctx).Table("x_user_enterprise").
		Select("e_id, COUNT(*) as count").
		Where("e_id IN (?)", eIds).
		Group("e_id").
		Scan(&userEnterpriseStats)

	eidToUserCount := make(map[uint32]int64)
	for _, stat := range userEnterpriseStats {
		eidToUserCount[stat.EID] = stat.Count
	}

	// 组装返回结果
	result := make([]*response.ListEnterpriseItem, 0, len(enterprises))
	for _, item := range enterprises {
		row := &response.ListEnterpriseItem{
			ID:        item.ID,
			FullName:  item.FullName,
			ShortName: item.ShortName,
			LeaseTime: item.LeaseTime.Format("2006-01-02"),
			Websites:  utils.Split(item.Website, ";"),
			Contact:   utils.Split(item.Contact, ";"),
			Status:    item.Status,
			Top:       item.Top,
			Weight:    item.Weight,
			Category:  []string{}, // 示例项目：类目字段留空
		}
		// 填入用户数信息
		if count, ok := eidToUserCount[item.ID]; ok {
			row.BNAME = fmt.Sprintf("用户数: %d", count) // 示例：借用 BNAME 字段显示用户数
		}
		result = append(result, row)
	}

	return result, total, nil
}

// UpdateEnterpriseStatus 更新企业状态
func (s *EnterpriseService) UpdateEnterpriseStatus(ctx context.Context, userId uint, params request.UpdateEnterpriseStatusParams) error {
	err := global.GVA_DB.WithContext(ctx).Model(&custom.XEnterprise{}).Where("id = ?", params.ID).Updates(map[string]any{
		"status":     params.Status,
		"updated_by": userId,
		"updated_at": time.Now(),
	}).Error
	if err != nil {
		return fmt.Errorf("更新状态失败，%v", err)
	}
	return nil
}

// TopEnterprise 企业置顶
func (s *EnterpriseService) TopEnterprise(ctx context.Context, userId uint, params request.TopEnterpriseParams) error {
	var topStatus int32
	if params.Status {
		topStatus = 1
	}
	err := global.GVA_DB.WithContext(ctx).Model(custom.XEnterprise{}).Where("id = ?", params.ID).Updates(map[string]any{
		"top":        topStatus,
		"updated_by": uint32(userId),
		"updated_at": time.Now(),
	}).Error
	if err != nil && params.Status {
		return fmt.Errorf("置顶失败，%v", err)
	}
	if err != nil && !params.Status {
		return fmt.Errorf("取消置顶失败，%v", err)
	}
	return nil
}
