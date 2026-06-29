package repository

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Cache), "*"),              // 缓存访问
	wire.Struct(new(EnterpriseRepo), "*"),     // 企业相关数据访问
	wire.Struct(new(UserRepo), "*"),           // 用户相关数据访问
	wire.Struct(new(UserEnterpriseRepo), "*"), // 用户关联的企业相关数据访问
)
