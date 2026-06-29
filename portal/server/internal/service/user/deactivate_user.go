package user

import (
	"context"
	"portal/internal/constant"
	"time"

	"portal/internal/model"
	"portal/internal/pkg"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/gormx"
	"portal/internal/pkg/log"
	"portal/internal/types"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// buildDeactivateUserCheck 构建查询用户注销准入评估结果
func buildDeactivateUserCheck(user *model.XUser) *types.DetailDeactivateUserCheckResp {
	canDeactivate := true
	items := make([]string, 0)

	return &types.DetailDeactivateUserCheckResp{
		CanDeactivate: canDeactivate,
		Items:         items,
	}
}

// DetailDeactivateUserCheck 注销前检查
func (s *Service) DetailDeactivateUserCheck(ctx context.Context) (*types.DetailDeactivateUserCheckResp, error) {
	user := pkg.DetailUserFromCtx(ctx)

	// 已注销或注销中，不允许重复注销
	if user.DeactivateStatus == constant.DeactivateUserStatusPending || user.DeactivateStatus == constant.DeactivateUserStatusFinished {
		return nil, errorx.UserAlreadyDeactivatedError
	}

	// 处理指标
	return buildDeactivateUserCheck(user), nil
}

// DeactivateUser 注销用户
// 先执行检查，若不通过则返回检查详情；若通过则更新状态并清理登录态
func (s *Service) DeactivateUser(ctx context.Context) (*types.DetailDeactivateUserCheckResp, error) {
	// 从上下文中获取当前用户信息
	user := pkg.DetailUserFromCtx(ctx)

	checkResp, err := s.DetailDeactivateUserCheck(ctx)
	if err != nil {
		return nil, err
	}

	// 检查未通过，直接返回检查结果
	if !checkResp.CanDeactivate {
		return checkResp, nil
	}

	// 所有指标通过，更新注销状态
	now := time.Now()
	err = s.UserRepo.MarkDeactivateStatus(ctx, user.ID, constant.DeactivateUserStatusPending, &now)
	if err != nil {
		log.WithContext(ctx).Error("更新用户注销状态失败", zap.Error(err))
		return nil, errorx.DeactivateUserdError
	}

	// 清理用户登录状态
	err = s.Cache.ClearLoginStatus(ctx, user.ID)
	if err != nil {
		log.WithContext(ctx).Error("清理用户登录状态失败", zap.Error(err))
	}
	return checkResp, nil
}

// DeactivateUserTaskHandler 用户注销处理任务
func (s *Service) DeactivateUserTaskHandler(ctx context.Context) {
	log.CronLoggerWithContext(ctx).Info("开始处理用户注销14天任务")

	// 分布式任务锁，防止多节点并发执行
	const taskKey = "user_deactivate_14d"
	// 预估执行时间：这里给 5 分钟的锁过期时间
	locked, err := s.Cache.AcquireTaskLock(ctx, taskKey, 5*time.Minute)
	if err != nil {
		log.CronLoggerWithContext(ctx).Error("获取用户注销14天任务分布式锁失败", zap.Error(err))
		return
	}
	if !locked {
		// 其他节点已经在执行本任务
		log.CronLoggerWithContext(ctx).Info("用户注销14天任务已在其他节点执行，当前节点跳过")
		return
	}
	defer func() {
		if err := s.Cache.ReleaseTaskLock(ctx, taskKey); err != nil {
			log.CronLoggerWithContext(ctx).Error("释放用户注销14天任务分布式锁失败", zap.Error(err))
		}
	}()

	// 计算 14 天之前的时间点
	before := time.Now().AddDate(0, 0, -14).Format("2006-01-02 15:04:05")

	// 查询 14 天前进入注销中的用户：deactivate_status = 1 且 deactivate_at <= before
	users, err := s.UserRepo.ListPendingDeactivateUsersBefore(ctx, before)
	if err != nil {
		log.CronLoggerWithContext(ctx).Error("查询待处理的注销用户失败", zap.Error(err))
		return
	}

	if len(users) == 0 {
		log.CronLoggerWithContext(ctx).Error("暂无需要处理的注销用户记录", zap.Error(err))
		return
	}

	for _, user := range users {
		// 每个用户一个事务：更新注销状态 + 删除实名认证记录
		err = gormx.Transaction(ctx, func(tx *gorm.DB) error {
			txCtx := gormx.NewContext(ctx, tx)

			// 更新注销信息和状态
			if err := s.UserRepo.UpdateDeactivateInfo(txCtx, user.ID); err != nil {
				return err
			}

			// 删除该用户在 x_user_identity_verification 中的记录
			if err := s.UserRepo.DeleteUserVerifyByUID(txCtx, user.ID); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.CronLoggerWithContext(ctx).Error("处理用户注销14天任务失败", zap.Uint32("user_id", user.ID), zap.Error(err))
			continue
		}
	}
}
