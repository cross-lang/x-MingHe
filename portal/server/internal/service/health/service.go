package health

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Service 健康检查服务
type Service struct {
	MysqlDB    *gorm.DB
	RedisDB    redis.UniversalClient
}

// HealthStatus 健康状态
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded HealthStatus = "degraded"
)

// ComponentStatus 组件状态
type ComponentStatus struct {
	Name      string       `json:"name"`
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message,omitempty"`
	LatencyMs int64        `json:"latency_ms"`
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string                     `json:"status"`
	Timestamp int64                      `json:"timestamp"`
	Components map[string]ComponentStatus `json:"components,omitempty"`
}

// CheckHealth 执行健康检查
func (s *Service) CheckHealth(ctx context.Context, detailed bool) (*HealthCheckResponse, error) {
	response := &HealthCheckResponse{
		Status:    string(HealthStatusHealthy),
		Timestamp: time.Now().Unix(),
	}

	if detailed {
		response.Components = make(map[string]ComponentStatus)
	}

	// 检查 MySQL 连接
	mysqlStatus := s.checkMySQL(ctx, detailed)
	if detailed {
		response.Components["mysql"] = mysqlStatus
	}
	if mysqlStatus.Status == HealthStatusUnhealthy {
		response.Status = string(HealthStatusUnhealthy)
	}

	// 检查 Redis 连接
	redisStatus := s.checkRedis(ctx, detailed)
	if detailed {
		response.Components["redis"] = redisStatus
	}
	if redisStatus.Status == HealthStatusUnhealthy {
		response.Status = string(HealthStatusUnhealthy)
	}

	// 如果有组件处于 degraded 状态且整体状态不是 unhealthy，则标记为 degraded
	if response.Status == string(HealthStatusHealthy) &&
		(mysqlStatus.Status == HealthStatusDegraded || redisStatus.Status == HealthStatusDegraded) {
		response.Status = string(HealthStatusDegraded)
	}

	return response, nil
}

// checkMySQL 检查 MySQL 连接
func (s *Service) checkMySQL(ctx context.Context, detailed bool) ComponentStatus {
	start := time.Now()
	status := ComponentStatus{
		Name:   "mysql",
		Status: HealthStatusHealthy,
	}

	// 使用 context 设置超时
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 获取底层 SQL DB
	sqlDB, err := s.MysqlDB.DB()
	if err != nil {
		status.Status = HealthStatusUnhealthy
		status.Message = fmt.Sprintf("获取数据库连接失败: %v", err)
		status.LatencyMs = time.Since(start).Milliseconds()
		return status
	}

	// 执行 ping 检查
	err = sqlDB.PingContext(ctx)
	if err != nil {
		status.Status = HealthStatusUnhealthy
		status.Message = fmt.Sprintf("数据库连接失败: %v", err)
		status.LatencyMs = time.Since(start).Milliseconds()
		return status
	}

	// 检查连接池状态
	if detailed {
		stats := sqlDB.Stats()
		status.Message = fmt.Sprintf(
			"连接池: 空闲=%d/最大=%d, 使用中=%d, 总数=%d",
			stats.Idle,
			stats.MaxOpenConnections,
			stats.InUse,
			stats.OpenConnections,
		)
	}

	status.LatencyMs = time.Since(start).Milliseconds()

	// 如果延迟超过 500ms，标记为 degraded
	if status.LatencyMs > 500 {
		status.Status = HealthStatusDegraded
		if detailed {
			status.Message = fmt.Sprintf("响应较慢: %s", status.Message)
		}
	}

	return status
}

// checkRedis 检查 Redis 连接
func (s *Service) checkRedis(ctx context.Context, detailed bool) ComponentStatus {
	start := time.Now()
	status := ComponentStatus{
		Name:   "redis",
		Status: HealthStatusHealthy,
	}

	// 使用 context 设置超时
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 执行 ping 检查
	err := s.RedisDB.Ping(ctx).Err()
	if err != nil {
		status.Status = HealthStatusUnhealthy
		status.Message = fmt.Sprintf("Redis 连接失败: %v", err)
		status.LatencyMs = time.Since(start).Milliseconds()
		return status
	}

	status.LatencyMs = time.Since(start).Milliseconds()

	// 如果延迟超过 100ms，标记为 degraded
	if status.LatencyMs > 100 {
		status.Status = HealthStatusDegraded
		if detailed {
			status.Message = "响应较慢"
		}
	}

	if detailed && status.Message == "" {
		poolStats := s.RedisDB.PoolStats()
		status.Message = fmt.Sprintf(
			"连接池: 命中=%d, 失败=%d, 总连接=%d, 空闲=%d",
			poolStats.Hits,
			poolStats.Misses,
			poolStats.TotalConns,
			poolStats.IdleConns,
		)
	}

	return status
}
