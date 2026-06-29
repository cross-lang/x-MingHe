package app

import (
	"portal/internal/config"
	enterpriseHandler "portal/internal/handler/enterprise"
	healthHandler "portal/internal/handler/health"
	userHandler "portal/internal/handler/user"
	"portal/internal/repository"
	enterpriseSve "portal/internal/service/enterprise"
	healthSve "portal/internal/service/health"
	userSve "portal/internal/service/user"

	"github.com/go-redis/redis_rate/v10"
	"gorm.io/gorm"
)

type DIContainer struct {
	// 配置
	Config config.Config

	// 处理器
	EnterpriseHandler *enterpriseHandler.Handler
	UserHandler       *userHandler.Handler
	HealthHandler     *healthHandler.Handler

	// 服务
	EnterpriseService *enterpriseSve.Service
	UserService       *userSve.Service
	HealthService     *healthSve.Service

	// 用户数据访问
	UserRepo       *repository.UserRepo
	EnterpriseRepo *repository.EnterpriseRepo

	// 缓存访问
	Cache *repository.Cache

	// 限流器
	RateLimiter *redis_rate.Limiter

	// MySQL
	MysqlDB *gorm.DB
}
