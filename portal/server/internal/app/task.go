package app

import (
	"context"
	"fmt"
	"portal/internal/pkg/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// StartTask 定时任务
func (a *App) StartTask() {
	// 注册用户注销14天任务，每天凌晨1点执行一次
	a.addTaskFunc("用户注销14天", "0 0 1 * * *", a.Container.UserService.DeactivateUserTaskHandler)
	// 启动定时任务
	a.Cron.Start()
}

// 添加任务处理方法和ctx状态
func (a *App) addTaskFunc(name string, spec string, cmd func(ctx context.Context)) {
	id, err := a.Cron.AddFunc(spec, func() {
		ctx := context.Background()
		ctx = log.ContextWithTraceID(ctx, uuid.New().String())
		cmd(ctx)
	})
	if err != nil {
		log.CronLogger.Error(fmt.Sprintf("[%s][%s]任务添加失败", name, spec), zap.Error(err))
	} else {
		log.CronLogger.Info(fmt.Sprintf("[%s][%s]任务添加成功", name, spec), zap.Int("id", int(id)))
	}
}
