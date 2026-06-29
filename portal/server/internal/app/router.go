package app

import (
	"portal/internal/middleware"
)

func (a *App) RegisterRoute() {
	// 静态路由
	a.Engine.Static("/page", "././view")

	// 应用信息路由
	a.Engine.GET("/info", InfoHandler)

	// 健康检查路由（无需认证）
	a.Engine.GET("/health", a.Container.HealthHandler.CheckHealth)
	a.Engine.GET("/ready", a.Container.HealthHandler.Readiness)
	a.Engine.GET("/live", a.Container.HealthHandler.Liveness)

	// 每分钟1次限流中间件
	// rateLimit1F1C := middleware.SpecializedRateLimitMiddleware(a.Container.RateLimiter, "1m", 1)

	// 用户登录身份检查中间件
	loginAuthMiddleware := middleware.LoginAuthMiddleware(
		a.Container.UserRepo,
		a.Container.Cache,
		a.Config.LoginJwt,
	)

	// 企业模块路由
	enterpriseGroup := a.Engine.Group("/v1/enterprises", loginAuthMiddleware)
	// 查询企业列表
	enterpriseGroup.GET("/list", a.Container.EnterpriseHandler.ListEnterprise)
	// 查询企业详情
	enterpriseGroup.GET("/:enterprise_id/detail", a.Container.EnterpriseHandler.DetailEnterprise)
	// 用户加入企业
	enterpriseGroup.POST("/join", a.Container.EnterpriseHandler.PostJoinEnterprise)

	// 用户模块路由
	userGroup := a.Engine.Group("/v1/users", loginAuthMiddleware)
	// 发送注册短信验证码
	a.Engine.POST("/v1/users/register/code/send", a.Container.UserHandler.PostSendRegisterCode)
	// 用户注册
	a.Engine.POST("/v1/users/register", a.Container.UserHandler.PostRegisterUser)
	// 发送登录短信验证码
	a.Engine.POST("/v1/users/login/code/send", a.Container.UserHandler.PostSendLoginCode)
	// 用户登录
	a.Engine.POST("/v1/users/login", a.Container.UserHandler.PostLoginUser)
	// 用户登出
	userGroup.POST("/logout", a.Container.UserHandler.PostLogoutUser)
	// 查询用户详情
	userGroup.GET("/detail", a.Container.UserHandler.DetailUser)
	// 查询用户入驻的企业列表
	userGroup.GET("/enterprises/list", a.Container.UserHandler.ListUserEnterprise)
	// 用户实名认证
	userGroup.POST("/verification", a.Container.UserHandler.PostVerifyUser)
	// 查询用户注销准入评估结果
	userGroup.GET("/check/deactivate/detail", a.Container.UserHandler.DetailCheckDeactivate)
	// 用户注销
	userGroup.POST("/deactivate", a.Container.UserHandler.PostDeactivateUser)
}
