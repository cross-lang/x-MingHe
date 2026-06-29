package repository

import (
	"context"
	"portal/internal/model"
	"portal/internal/pkg/gormx"

	"gorm.io/gorm"
)

// EnterpriseRepo 企业访问层
type EnterpriseRepo struct {
	DB *gorm.DB
}

// SearchEnterpriseListParams 企业列表搜索参数
type SearchEnterpriseListParams struct {
	ParkId   uint32 // 园区ID
	Type     int    // 企业类型（1：企业，2：学院）
	Keyword  string // 搜索关键词
	Category int    // 类目ID
	Size     int    // 每页数量
	Page     int    // 页码
}

// SearchEnterpriseList 获取指定园区下的企业列表
func (r *EnterpriseRepo) SearchEnterpriseList(ctx context.Context, searchParams *SearchEnterpriseListParams) ([]*model.XEnterprise, int64, error) {
	query := r.DB.Table("x_enterprise AS a").Select("a.*")
	query = query.Joins("LEFT JOIN x_enterprise_category_relation AS c ON c.e_id = a.id")
	query = query.Where("a.p_id = ?", searchParams.ParkId)
	query = query.Where("a.type = ?", searchParams.Type)

	if searchParams.Keyword != "" {
		keyword := "%" + searchParams.Keyword + "%"
		query = query.Where("a.full_name LIKE ? OR a.short_name LIKE ? OR a.label LIKE ?", keyword, keyword, keyword)
	}

	if searchParams.Category != 0 {
		query = query.Where("c.c_id = ?", searchParams.Category)
	}

	// 软删除过滤
	query = query.Where("a.deleted_at IS NULL")

	// 查询总数
	var count int64
	err := r.DB.WithContext(ctx).Table("(?) as data", query).Group("data.id").Select("COUNT(data.id)").Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询分页列表
	offset, size := gormx.Pagination(searchParams.Size, searchParams.Page)
	order := "data.top desc, data.weight desc, data.updated_at desc"
	var enterprises []*model.XEnterprise
	err = r.DB.WithContext(ctx).Table("(?) AS data", query).Group("data.id").Order(order).Offset(offset).Limit(size).Find(&enterprises).Error
	return enterprises, count, err
}

// GetById 通过ID获取企业信息
func (r *EnterpriseRepo) GetById(ctx context.Context, id uint32) (*model.XEnterprise, error) {
	var enterprise model.XEnterprise
	err := r.DB.WithContext(ctx).First(&enterprise, "id = ?", id).Error
	return &enterprise, err
}
