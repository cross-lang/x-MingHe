package health

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet Health service provider set
var ProviderSet = wire.NewSet(
	NewService,
)

// NewService 创建健康检查服务
func NewService(db *gorm.DB, rdb redis.UniversalClient) *Service {
	return &Service{
		MysqlDB: db,
		RedisDB: rdb,
	}
}
