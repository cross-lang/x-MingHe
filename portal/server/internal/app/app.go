package app

import (
	"context"
	"errors"
	"fmt"
	logg "log"
	"net/http"
	"os"
	"os/signal"
	"portal/internal/config"
	"portal/internal/middleware"
	"portal/internal/pkg/log"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	Config    config.Config
	Engine    *gin.Engine
	Serve     *http.Server
	Container *DIContainer // DI容器
	Cron      *cron.Cron
}

func InitApp(diContainer *DIContainer) *App {
	App := &App{
		Config:    diContainer.Config,
		Engine:    gin.New(),
		Container: diContainer,
		Cron:      cron.New(cron.WithSeconds()),
	}
	conf := diContainer.Config

	// 根据调试模式设置 Gin 模式
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化日志
	log.InitLogger(conf.Logger)

	// 初始化web服务中间件
	App.Engine.Use(middleware.Recovery())
	App.Engine.Use(middleware.CorsAll())
	App.Engine.Use(middleware.TraceMiddleware())
	App.Engine.Use(middleware.DebugMiddleware(conf))

	// 调试模式下不启用限流
	if !conf.Debug {
		App.Engine.Use(middleware.GlobalRateLimitMiddleware(diContainer.RateLimiter, conf.GlobalRateLimit))
	}

	// 调试模式下记录请求和响应 body
	if conf.Debug {
		// 添加所有内容类型到收集列表
		allContentTypes := []string{"application/json", "application/x-www-form-urlencoded", "multipart/form-data", "text/plain"}
		App.Engine.Use(middleware.LoggerMiddleware(allContentTypes))
	} else {
		App.Engine.Use(middleware.LoggerMiddleware(conf.CollectBodyContentType))
	}

	// 注册Swagger路由
	App.RegisterSwagger()

	// 注册路由
	App.RegisterRoute()

	// 迁移表结构
	App.AutoMigrate()

	// 启动定时任务
	App.StartTask()

	return App
}

func (a *App) RunPostDeactivateUserServer() {
	a.Serve = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.ServerPort),
		Handler: a.Engine.Handler(),
	}

	if err := a.Serve.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Fatalf("listen: %s\n", err)
	}
}

func (a *App) Run() {
	go a.RunPostDeactivateUserServer()

	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.Serve.Shutdown(ctx); err != nil {
		logg.Println("Server Shutdown:", err)
	}
	logg.Println("Server exiting")
}
