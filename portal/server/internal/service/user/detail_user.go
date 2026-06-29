package user

import (
	"context"
	"portal/internal/constant"
	"portal/internal/pkg"
	"portal/internal/pkg/stringx"
	"portal/internal/types"
)

// DetailUser 查询用户详情
func (s *Service) DetailUser(ctx context.Context) (*types.UserItem, error) {
	// 从上下文中查询用户详情
	user := pkg.DetailUserFromCtx(ctx)

	// 获取用户当前所在的企业和园区
	space, err := s.UserRepo.GetUserSpace(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// 计算用户角色
	roles := make([]*types.UserRole, 0)

	// 如果用户没有加入企业并且未完成实名认证
	if space["e_id"] == 0 && user.VerificationStatus != 1 {
		// 访客
		roles = append(roles, &types.UserRole{Key: constant.TouristRole, Name: "访客"})
	}
	// 如果用户没有加入企业并且已完成实名认证
	if space["e_id"] == 0 && user.VerificationStatus == 1 {
		// 访客
		roles = append(roles, &types.UserRole{Key: constant.VisitorRole, Name: "游客"})
	}

	// 如果用户已经加入企业查询用户在企业的角色
	if space["e_id"] != 0 {
		roles, err = s.getUserEnterpriseRoles(ctx, user.ID, space["e_id"])
		if err != nil {
			return nil, err
		}
	}

	result := &types.UserItem{
		ID:                  user.ID,
		Name:                user.Name,
		Avatar:              user.Avatar,
		Account:             user.Account,
		PhoneNumber:         user.PhoneNumber,
		Gender:              user.Gender,
		VerificationStatus:  user.VerificationStatus,
		AccountStatus:       user.AccountStatus,
		BlockStatus:         user.BlockStatus,
		Roles:               roles,
		PID:                 space["p_id"],
		EID:                 space["e_id"],
		Profession:          user.Profession,
		GraduatedSchool:     user.GraduatedSchool,
		WorkYears:           int32(user.WorkYears),
		Residence:           user.Residence,
		Skills:              stringx.Split(user.Skills, ";"),
		Introduction:        user.Introduction,
		Contact:             user.Contact,
		EducationExperience: make([]*types.UserExperience, 0),
		WorkExperience:      make([]*types.UserExperience, 0),
		ProjectExperience:   make([]*types.UserExperience, 0),
	}

	return result, nil
}

// getUserEnterpriseRoles 获取用户在企业中的角色
func (s *Service) getUserEnterpriseRoles(ctx context.Context, userId, enterpriseId uint32) ([]*types.UserRole, error) {
	// TODO: 实现获取用户角色的逻辑
	// 暂时返回默认角色
	return []*types.UserRole{
		{Key: "enterprise_staff", Name: "企业员工"},
	}, nil
}
