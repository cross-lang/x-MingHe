package repository

import (
	"context"
	"portal/internal/model"
	"portal/internal/pkg/gormx"

	"gorm.io/gorm"
)

// UserEnterpriseRepo 用户和企业关联访问层
type UserEnterpriseRepo struct {
	DB *gorm.DB
}

// GetUserEnterprises 查询用户入驻的企业列表
func (r *UserEnterpriseRepo) GetUserEnterprises(ctx context.Context, userId uint32) ([]*model.XEnterprise, error) {
	query := r.DB.WithContext(ctx).Table("x_enterprise AS a")
	query = query.Joins("LEFT JOIN x_user_enterprise AS b ON b.e_id = a.id AND b.deleted_at IS NULL")
	query = query.Where("b.u_id = ?", userId)
	var enterprises []*model.XEnterprise
	err := query.Find(&enterprises).Error
	return enterprises, err
}

// SaveUserEnterprise 保存用户入职企业或者学院关系
func (r *UserEnterpriseRepo) SaveUserEnterprise(ctx context.Context, userEnterprise *model.XUserEnterprise) error {
	db := gormx.WithContext(ctx, r.DB)
	return db.WithContext(ctx).Save(userEnterprise).Error
}

// UserHasEnterprise 检查用户是否已经入职了企业
func (r *UserEnterpriseRepo) UserHasEnterprise(ctx context.Context, userId uint32, enterpriseId uint32) (bool, error) {
	var count int64
	query := r.DB.WithContext(ctx).Table(model.TableNameXUserEnterprise)
	query = query.Where("u_id = ? AND e_id = ?", userId, enterpriseId)
	query = query.Where("deleted_at IS NULL")
	err := query.Count(&count).Error
	return count > 0, err
}
