package handler

import (
	"portal/internal/handler/enterprise"
	"portal/internal/handler/health"
	"portal/internal/handler/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(enterprise.Handler), "*"), // 企业相关接口处理器
	wire.Struct(new(user.Handler), "*"),       // 用户相关接口处理器
	health.ProviderSet,                        // 健康检查处理器
)
