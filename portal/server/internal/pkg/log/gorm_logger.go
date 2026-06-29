package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger GORM 日志适配器
type GormLogger struct {
	Debug bool // 是否启用调试模式
}

// NewGormLogger 创建 GORM 日志器
func NewGormLogger(debug bool) *GormLogger {
	return &GormLogger{Debug: debug}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	WithContext(ctx).Info(fmt.Sprintf(s, i...))
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	WithContext(ctx).Warn(fmt.Sprintf(s, i...))
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	WithContext(ctx).Error(fmt.Sprintf(s, i...))
}

// Trace 追踪 SQL 执行
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 调试模式下记录所有 SQL 查询
	if l.Debug {
		logFields := []zap.Field{
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("duration", elapsed),
			zap.String("module", string(LogModuleDatabase)),
		}

		if err != nil {
			logFields = append(logFields, zap.Error(err))
			WithModule(ctx, LogModuleDatabase).Error("SQL执行失败", logFields...)
		} else {
			WithModule(ctx, LogModuleDatabase).Debug("SQL执行成功", logFields...)
		}
		return
	}

	// 非调试模式，只记录错误和慢查询
	if err != nil {
		WithContext(ctx).Error("SQL执行失败",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("time", int64(elapsed.Milliseconds())),
			zap.String("module", string(LogModuleDatabase)),
		)
	} else if elapsed > 200*time.Millisecond {
		// 慢查询警告
		WithContext(ctx).Warn("慢SQL查询",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Int64("time", int64(elapsed.Milliseconds())),
			zap.String("module", string(LogModuleDatabase)),
		)
	}
}

