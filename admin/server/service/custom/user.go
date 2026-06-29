package custom

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/constant"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/encrypt"
	"gorm.io/gorm"
)

type UserService struct {
}

// UpdateUser 编辑用户
func (s *UserService) UpdateUser(ctx context.Context, req *request.UpdateUserReq) error {
	// 查询用户
	var user custom.XUser
	if err := global.GVA_DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", req.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}
	if user.DeactivateStatus != 0 {
		return errors.New("user is deactivated")
	}

	err := global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 更新用户信息和状态
		err := tx.Model(&custom.XUser{}).Where("id = ?", req.ID).Updates(map[string]any{
			"account":          req.PhoneNumber,
			"phone_number":     req.PhoneNumber,
			"avatar":           req.Avatar,
			"graduated_school": req.GraduatedSchool,
			"profession":       req.Profession,
			"residence":        req.Residence,
			"contact":          req.Contact,
			"work_years":       req.WorkYears,
			"introduction":     req.Introduction,
			"skills":           strings.Join(req.Skills, ";"),
			"block_status":     req.BlockStatus,
			"account_status":   req.AccountStatus,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}


// DisableUser 禁用用户
func (s *UserService) DisableUser(ctx context.Context, req *request.DisableUserReq) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		db := global.GVA_DB.WithContext(ctx)
		if err := db.Model(&custom.XUser{}).Where("id = ?", req.UserID).Update("account_status", constant.UserAccountStatusDisable).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		return nil
	})
}

// EnableUser 启用用户
func (s *UserService) EnableUser(ctx context.Context, req *request.EnableUserReq) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		db := global.GVA_DB.WithContext(ctx)
		if err := db.Model(&custom.XUser{}).Where("id = ?", req.UserID).Update("account_status", constant.UserAccountStatusNormal).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		return nil
	})
}

// DetailUser 查询用户详情
func (s *UserService) DetailUser(ctx context.Context, userId uint, req *request.DetailUserReq) (*response.UserItem, error) {
	db := global.GVA_DB.WithContext(ctx)

	// 查询用户
	var user custom.XUser
	if err := db.Where("id = ? AND deleted_at IS NULL", req.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// 查询用户实名认证
	var userVerify custom.XUserIdentityVerification
	if err := db.Where("u_id = ? AND deleted_at IS NULL", user.ID).First(&userVerify).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	var deactivateAt = ""
	if user.DeactivateAt != nil {
		deactivateAt = user.DeactivateAt.Format("2006-01-02 15:04:05")
	}
	// 组装最终结果
	result := &response.UserItem{
		ID:                 user.ID,
		Name:               user.Name,
		Avatar:             user.Avatar,
		Account:            user.Account,
		PhoneNumber:        user.PhoneNumber,
		Gender:             user.Gender,
		VerificationStatus: user.VerificationStatus,
		AccountStatus:      user.AccountStatus,
		BlockStatus:        user.BlockStatus,
		DeactivateStatus:   user.DeactivateStatus,
		DeactivateAt:       deactivateAt,
		CreatedAt:          user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          user.UpdatedAt.Format("2006-01-02 15:04:05"),

		VChannel: "",
		VName:    "",
		VIDCard:  "",

		Profession:          user.Profession,
		GraduatedSchool:     user.GraduatedSchool,
		WorkYears:           user.WorkYears,
		Residence:           user.Residence,
		Skills:              utils.Split(user.Skills, ";"),
		Introduction:        user.Introduction,
		Contact:             user.Contact,
	}

	// 用户有注销状态时
	if user.DeactivateStatus != 0 {
		re := regexp.MustCompile(`\([^)]*\)`)
		result.Account = re.ReplaceAllString(user.Account, "")
		result.PhoneNumber = re.ReplaceAllString(user.PhoneNumber, "")
		result.DeactivateAt = user.DeactivateAt.Format("2006-01-02 15:04:05")
	}

	// 填充实名认证信息
	if userVerify.ID > 0 {
		result.VChannel = userVerify.Channel
		result.VName = userVerify.Name
		result.VIDCard = userVerify.IDCard

		// 解密身份证号码
		iDCard, err := encrypt.AESDecrypt(userVerify.IDCard, []byte(global.GVA_CONFIG.System.DataEncryptKey))
		if err != nil {
			return nil, err
		}
		// 如果是超级管理员展示全部的
		isSuperAdmin := utils.IsSuperAdmin(ctx, userId)
		if isSuperAdmin {
			result.VIDCard = string(iDCard)
		} else {
			if len(iDCard) == 18 {
				result.VIDCard = string(iDCard[:6]) + "********" + string(iDCard[14:])
			}
		}
	}

	return result, nil
}

// ListUser 查询用户列表
func (s *UserService) ListUser(ctx context.Context, userId uint, req *request.ListUserReq) (*response.UserList, error) {
	db := global.GVA_DB.WithContext(ctx)

	// Step 1: 构建基础查询
	query := db.Table("x_user as u").Where("u.deleted_at IS NULL")

	// Step 2.1: 账号状态过滤
	if len(req.AccountStatuses) > 0 {
		query = query.Where("u.account_status IN ?", req.AccountStatuses)
	}

	// Step 2.2: 拉黑状态过滤
	if len(req.BlockStatuses) > 0 {
		query = query.Where("u.block_status IN ?", req.BlockStatuses)
	}

	// Step 2.3: 实名状态过滤
	if len(req.VerifyStatuses) > 0 {
		query = query.Where("u.verification_status IN ?", req.VerifyStatuses)
	}

	if len(req.DeactivateStatus) > 0 {
		query = query.Where("u.deactivate_status IN ?", req.DeactivateStatus)
	}

	// Step 3: 注册时间范围过滤
	if req.RegisterTimeFrom != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.RegisterTimeFrom); err == nil {
			query = query.Where("u.created_at >= ?", t)
		}
	}
	if req.RegisterTimeTo != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.RegisterTimeTo); err == nil {
			query = query.Where("u.created_at <= ?", t)
		}
	}

	// Step 4: 关键字搜索（匹配账号、手机号、用户名）
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("u.account LIKE ? OR u.phone_number LIKE ? OR u.name LIKE ?", keyword, keyword, keyword)
	}

	// Step 5: 排序
	orderClause := "u.created_at DESC"
	if req.OrderKey != "" {
		switch req.OrderKey {
		case "id":
			orderClause = "u.id"
		case "name":
			orderClause = "u.name"
		case "account":
			orderClause = "u.account"
		case "phone_number":
			orderClause = "u.phone_number"
		case "created_at":
			orderClause = "u.created_at"
		default:
			orderClause = "u.created_at"
		}
		if req.Desc {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}
	}
	query = query.Order(orderClause)

	// Step 6: 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Step 7: 分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// Step 8: 查询用户基本信息
	var users []struct {
		ID                 uint32
		Name               string
		Avatar             string
		Account            string
		PhoneNumber        string
		Gender             int32
		VerificationStatus int32
		AccountStatus      int32
		BlockStatus        int32
		DeactivateStatus   int32
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}
	if err := query.Select([]string{
		"u.id",
		"u.name",
		"u.avatar",
		"u.account",
		"u.phone_number",
		"u.gender",
		"u.verification_status",
		"u.account_status",
		"u.block_status",
		"u.deactivate_status",
		"u.created_at",
		"u.updated_at",
	}).Offset(offset).Limit(pageSize).Scan(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return &response.UserList{Count: total, Items: []*response.UserItem{}}, nil
	}

	// Step 9: 批量收集用户ID
	userIDs := make([]uint32, 0, len(users))
	userMap := make(map[uint32]*response.UserItem)
	re := regexp.MustCompile(`\([^)]*\)`)
	for _, u := range users {
		account := re.ReplaceAllString(u.Account, "")
		phoneNumber := re.ReplaceAllString(u.PhoneNumber, "")
		item := &response.UserItem{
			ID:                 u.ID,
			Name:               u.Name,
			Avatar:             u.Avatar,
			Account:            account,
			PhoneNumber:        phoneNumber,
			Gender:             u.Gender,
			VerificationStatus: u.VerificationStatus,
			AccountStatus:      u.AccountStatus,
			BlockStatus:        u.BlockStatus,
			DeactivateStatus:   u.DeactivateStatus,
			CreatedAt:          u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:          u.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		userIDs = append(userIDs, u.ID)
		userMap[u.ID] = item
	}

	// Step 10: 批量查询实名认证信息
	var verifications []struct {
		UID     uint32 `gorm:"column:u_id"`
		Channel string `gorm:"column:channel"`
		Name    string `gorm:"column:name"`
		IDCard  string `gorm:"column:id_card"`
	}
	if err := db.Table("x_user_identity_verification").
		Where("u_id IN ? AND deleted_at IS NULL", userIDs).
		Order("created_at DESC").
		Find(&verifications).Error; err != nil {
		return nil, err
	}

	verifyMap := make(map[uint32]struct {
		Channel string
		Name    string
		IDCard  string
	})
	for _, v := range verifications {
		if _, exists := verifyMap[v.UID]; !exists {
			verifyMap[v.UID] = struct {
				Channel string
				Name    string
				IDCard  string
			}{v.Channel, v.Name, v.IDCard}
		}
	}

	// 查询用户是否是超级管理员
	isSuperAdmin := utils.IsSuperAdmin(ctx, userId)
	for uid, info := range verifyMap {
		if user, ok := userMap[uid]; ok {
			user.VChannel = info.Channel
			user.VName = info.Name
			iDCard, err := encrypt.AESDecrypt(info.IDCard, []byte(global.GVA_CONFIG.System.DataEncryptKey))
			if err != nil {
				continue
			}
			// 如果是超级管理员展示全部的
			if isSuperAdmin {
				user.VIDCard = string(iDCard)
			} else {
				if len(iDCard) == 18 {
					user.VIDCard = string(iDCard[:6]) + "********" + string(iDCard[14:])
				}
			}
		}
	}


	// Step 11: 组装最终结果
	items := make([]*response.UserItem, 0, len(users))
	for _, u := range users {
		if item, ok := userMap[u.ID]; ok {
			items = append(items, item)
		}
	}

	return &response.UserList{
		Count:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Items:    items,
	}, nil
}
