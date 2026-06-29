package service

import (
	"portal/internal/service/enterprise"
	"portal/internal/service/health"
	"portal/internal/service/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(enterprise.Service), "*"), // 企业业务层
	wire.Struct(new(user.Service), "*"),       // 用户业务层
	health.ProviderSet,                        // 健康检查服务
)
