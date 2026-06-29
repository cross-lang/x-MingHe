package log

// LogLevel 日志级别类型
type LogLevel string

const (
	// LogLevelDebug 调试级别，用于详细的诊断信息
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo 信息级别，用于一般性信息消息
	LogLevelInfo LogLevel = "info"
	// LogLevelWarning 警告级别，用于潜在的有害情况
	LogLevelWarning LogLevel = "warn"
	// LogLevelError 错误级别，用于错误事件，但应用仍可继续运行
	LogLevelError LogLevel = "error"
	// LogLevelCritical 严重级别，用于可能导致应用停止的严重事件
	LogLevelCritical LogLevel = "critical"
	// LogLevelPanic 恐慌级别，用于致命错误后应用退出
	LogLevelPanic LogLevel = "panic"
	// LogLevelFatal 致命级别，用于致命错误后应用退出
	LogLevelFatal LogLevel = "fatal"
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	return string(l)
}

// IsValid 检查日志级别是否有效
func (l LogLevel) IsValid() bool {
	switch l {
	case LogLevelDebug, LogLevelInfo, LogLevelWarning,
		LogLevelError, LogLevelCritical, LogLevelPanic, LogLevelFatal:
		return true
	default:
		return false
	}
}

// ParseLevel 从字符串解析日志级别
func ParseLevel(level string) LogLevel {
	switch level {
	case "debug", "DEBUG":
		return LogLevelDebug
	case "info", "INFO":
		return LogLevelInfo
	case "warn", "WARNING", "warning":
		return LogLevelWarning
	case "error", "ERROR":
		return LogLevelError
	case "critical", "CRITICAL":
		return LogLevelCritical
	case "panic", "PANIC":
		return LogLevelPanic
	case "fatal", "FATAL":
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}

// LogModule 日志模块类型
type LogModule string

const (
	// LogModuleAPI API 模块
	LogModuleAPI LogModule = "api"
	// LogModuleService 服务模块
	LogModuleService LogModule = "service"
	// LogModuleRepository 数据访问模块
	LogModuleRepository LogModule = "repository"
	// LogModuleMiddleware 中间件模块
	LogModuleMiddleware LogModule = "middleware"
	// LogModuleConfig 配置模块
	LogModuleConfig LogModule = "config"
	// LogModuleDatabase 数据库模块
	LogModuleDatabase LogModule = "database"
	// LogModuleCache 缓存模块
	LogModuleCache LogModule = "cache"
	// LogModuleExternal 外部服务模块
	LogModuleExternal LogModule = "external"
	// LogModuleSystem 系统模块
	LogModuleSystem LogModule = "system"
)

// String 返回模块的字符串表示
func (m LogModule) String() string {
	return string(m)
}
