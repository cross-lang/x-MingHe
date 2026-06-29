package repository

import (
	"context"
	"errors"
	"fmt"
	"portal/internal/constant"
	"portal/internal/model"
	"portal/internal/pkg/gormx"
	"time"

	"gorm.io/gorm"
)

// UserRepo 用户访问层
type UserRepo struct {
	DB *gorm.DB
}

// GetById 通过ID获取用户
func (r *UserRepo) GetById(ctx context.Context, id uint32) (*model.XUser, error) {
	var user model.XUser
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

// GetByPhoneNumber 通过手机号获取用户信息
func (r *UserRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.XUser, error) {
	var user model.XUser
	err := r.DB.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user).Error
	return &user, err
}

// PhoneNumberExists 检查手机号是否已经注册
func (r *UserRepo) PhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Table(model.TableNameXUser).Where("phone_number = ?", phoneNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RegisterSave 保存注册用户信息
func (r *UserRepo) RegisterSave(ctx context.Context, data *model.XUser) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

// UpdateVerifyStatus 更新用户实名认证状态
func (r *UserRepo) UpdateVerifyStatus(ctx context.Context, id uint32, name string) error {
	db := r.DB.WithContext(ctx)
	return db.Model(model.XUser{}).Where("id = ?", id).Updates(map[string]any{
		"verify": 1,    // 认证状态
		"name":   name, // 认证后的昵称
	}).Error
}

// SaveVerify 保存用户实名认证信息
func (r *UserRepo) SaveVerify(ctx context.Context, data *model.XUserIdentityVerification) error {
	db := r.DB.WithContext(ctx)
	return db.Create(data).Error
}

// IDCardHasVerify 检查身份证号码是否已认证
func (r *UserRepo) IDCardHasVerify(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(model.XUserIdentityVerification{}).Where("hash = ?", hash).Count(&count).Error
	return count > 0, err
}

// GetUserSpace 获取用户当前所处园区和企业空间状态
func (r *UserRepo) GetUserSpace(ctx context.Context, userId uint32) (map[string]uint32, error) {
	var result struct {
		EID uint32
		PID uint32
	}
	err := r.DB.WithContext(ctx).Table("x_user_enterprise").Where("u_id = ?", userId).Select("e_id, p_id").First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]uint32{"e_id": 0, "p_id": 0}, nil
		}
		return nil, err
	}
	return map[string]uint32{"e_id": result.EID, "p_id": result.PID}, nil
}

// MarkDeactivateStatus 标记用户注销状态
func (r *UserRepo) MarkDeactivateStatus(ctx context.Context, userId uint32, status int32, deactivateAt *time.Time) error {
	db := gormx.WithContext(ctx, r.DB)
	data := map[string]any{
		"deactivate_status": status,
		"deactivate_at":     deactivateAt,
	}
	return db.Model(model.XUser{}).Where("id = ?", userId).Updates(data).Error
}

// UpdateDeactivateInfo 更新用户注销信息
func (r *UserRepo) UpdateDeactivateInfo(ctx context.Context, userId uint32) error {
	db := gormx.WithContext(ctx, r.DB)
	var user model.XUser
	err := db.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return err
	}

	now := time.Now()
	data := map[string]any{
		"deactivate_status": constant.DeactivateUserStatusFinished, // 已注销
		"deactivate_at":     &now,                                  // 注销时间
		"account":           user.Account + fmt.Sprintf("(已注销%s)", now.Format("2006-01-02 15:04:05")),
		"phone_number":      user.PhoneNumber + fmt.Sprintf("(已注销%s)", now.Format("2006-01-02 15:04:05")),
	}
	return db.Model(model.XUser{}).Where("id = ?", userId).Updates(data).Error
}

// DeleteUserVerifyByUID 根据用户ID删除实名认证记录
func (r *UserRepo) DeleteUserVerifyByUID(ctx context.Context, userId uint32) error {
	db := gormx.WithContext(ctx, r.DB)
	return db.WithContext(ctx).Where("u_id = ?", userId).Delete(&model.XUserIdentityVerification{}).Error
}

// ListPendingDeactivateUsersBefore 查询在指定时间点之前进入注销中的用户
func (r *UserRepo) ListPendingDeactivateUsersBefore(ctx context.Context, before string) ([]*model.XUser, error) {
	var users []*model.XUser
	err := r.DB.WithContext(ctx).Model(&model.XUser{}).Where("deactivate_status = ? AND deactivate_at <= ?", 1, before).Find(&users).Error
	return users, err
}
