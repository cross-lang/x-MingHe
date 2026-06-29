package utils

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// IsSuperAdmin 判断用户身份是超级管理员
func IsSuperAdmin(ctx context.Context, userId uint) bool {
	query := global.GVA_DB.WithContext(ctx).Model(&system.SysUserAuthority{})
	query = query.Where("sys_user_id = ?", userId)
	query = query.Where("sys_authority_authority_id = ?", "888")
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
