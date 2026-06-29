//go:build wireinject
// +build wireinject

package main

import (
	"portal/internal/app"
	"portal/internal/config"
	"portal/internal/handler"
	"portal/internal/infra"
	"portal/internal/repository"
	"portal/internal/service"

	"github.com/google/wire"
)

// 配置项注入
var configSet = wire.NewSet(
	config.LoadConfig,      // 加载config
	config.LoadMysqlConfig, // 加载mysql配置
	config.LoadRedisConfig, // 加载redis配置
)

// 存储连接和相关工具注入
var dbSet = wire.NewSet(
	infra.InitMysql,        // 初始化mysql连接
	infra.InitRedis,        // 初始化redis连接
	infra.InitRedisLimiter, // 初始化限流器
)

func wireDIContainer() (*app.DIContainer, error) {

	panic(wire.Build(
		// 配置项注入
		configSet,
		// 存储连接和相关工具注入
		dbSet,
		// 数据访问层注入
		repository.ProviderSet,
		// 业务逻辑层注入
		service.ProviderSet,
		// API 接口层注入
		handler.ProviderSet,
		// App 全局di容器注入
		wire.NewSet(wire.Struct(new(app.DIContainer), "*")),
	))
}
