package log

import (
	"context"
	"io"
	"log"
	"os"
	"portal/internal/config"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var CronLogger *zap.Logger

type RemoteCore struct {
	core         zapcore.Core
	enableRemote bool
	remoteURL    string
}

func (r *RemoteCore) Enabled(level zapcore.Level) bool {
	return r.core.Enabled(level)
}

func (r *RemoteCore) With(fields []zapcore.Field) zapcore.Core {
	return &RemoteCore{
		core:         r.core.With(fields),
		enableRemote: r.enableRemote,
		remoteURL:    r.remoteURL,
	}
}

func (r *RemoteCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if r.Enabled(ent.Level) {
		return ce.AddCore(ent, r)
	}
	return ce
}

func (r *RemoteCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	err := r.core.Write(entry, fields)
	if err != nil {
		return err
	}
	// 异步推送到远程日志服务
	if r.enableRemote && r.remoteURL != "" {
		go r.pushRemote(entry, fields)
	}
	return nil
}

func (r *RemoteCore) Sync() error {
	return r.core.Sync()
}

// 推送日志到远程
func (r *RemoteCore) pushRemote(entry zapcore.Entry, fields []zapcore.Field) {
	log.Println("推送日志服务")
	//req, _ := http.NewRequest("POST", r.remoteURL, strings.NewReader(entry.Message))
	//req.Header.Set("Content-Type", "PostDeactivateUserlication/json")
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
}

var levelSet = map[string]zapcore.Level{
	"debug":    zapcore.DebugLevel,
	"info":     zapcore.InfoLevel,
	"warn":     zapcore.WarnLevel,
	"warning":  zapcore.WarnLevel,
	"error":    zapcore.ErrorLevel,
	"critical": zapcore.ErrorLevel, // critical 映射到 error level
	"panic":    zapcore.PanicLevel,
	"fatal":    zapcore.FatalLevel,
}

// InitLogger 初始化日志模块
func InitLogger(conf config.LoggerConfig) {
	// 确保日志目录存在
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("创建日志目录失败: %v", err)
	}

	// 项目日志文件名: MingHe_portal_<年月日时分秒毫秒>.log
	projectLogFile := "logs/MingHe_portal_" + time.Now().Format("20060102150405.000") + ".log"
	lumberjackLogger := &lumberjack.Logger{
		Filename:   projectLogFile,
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     7, // 天
		Compress:   true,
	}
	writeSyncer := zapcore.AddSync(io.MultiWriter(lumberjackLogger, os.Stdout))

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000")) // 精确到毫秒
	}

	level, ok := levelSet[conf.Level]
	if !ok {
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)
	// 用 remoteCore 包装原 core
	Logger = zap.New(zapcore.NewTee(
		&RemoteCore{
			core:         core,
			enableRemote: conf.EnableRemote,
			remoteURL:    conf.RemoteURL,
		}), zap.AddCaller())

	// 初始化定时任务日志器
	initCronLogger(conf)
}

// InitCronLogger 初始化定时任务日志器
func initCronLogger(conf config.LoggerConfig) {
	// 定时任务日志文件名: MingHe_portal_cron_<年月日时分秒毫秒>.log
	cronLogFile := "logs/MingHe_portal_cron_" + time.Now().Format("20060102150405.000") + ".log"
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cronLogFile,
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     7, // 天
		Compress:   true,
	}
	writeSyncer := zapcore.AddSync(io.MultiWriter(lumberjackLogger, os.Stdout))

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000")) // 精确到毫秒
	}

	level, ok := levelSet[conf.Level]
	if !ok {
		level = zap.InfoLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)
	CronLogger = zap.New(core, zap.AddCaller())
}

const TraceIDKey = "trace_id"

// TraceIDFromContext 从 Context 获取 trace_id
func TraceIDFromContext(ctx context.Context) string {
	if v := ctx.Value(TraceIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// ContextWithTraceID 在 Context 中设置 trace_id
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithContext 从 Context 获取 trace_id 并返回日志器
func WithContext(ctx context.Context) *zap.Logger {
	return WithModule(ctx, "")
}

// WithModule 从 Context 获取 trace_id 并添加模块信息，返回日志器
func WithModule(ctx context.Context, module LogModule) *zap.Logger {
	logger := Logger
	traceID := TraceIDFromContext(ctx)
	if traceID != "" {
		logger = logger.With(zap.String("trace_id", traceID))
	}
	if module != "" {
		logger = logger.With(zap.String("module", string(module)))
	}
	return logger
}

// WithFields 添加多个字段到日志器
func WithFields(ctx context.Context, module LogModule, fields ...zap.Field) *zap.Logger {
	logger := WithModule(ctx, module)
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	return logger
}

// CronLoggerWithContext 从 Context 获取 trace_id 并返回定时任务日志器
func CronLoggerWithContext(ctx context.Context) *zap.Logger {
	return CronLoggerWithModule(ctx, "")
}

// CronLoggerWithModule 从 Context 获取 trace_id 并添加模块信息，返回定时任务日志器
func CronLoggerWithModule(ctx context.Context, module LogModule) *zap.Logger {
	logger := CronLogger
	traceID := TraceIDFromContext(ctx)
	if traceID != "" {
		logger = logger.With(zap.String("trace_id", traceID))
	}
	if module != "" {
		logger = logger.With(zap.String("module", string(module)))
	}
	return logger
}

// CronLoggerWithFields 添加多个字段到定时任务日志器
func CronLoggerWithFields(ctx context.Context, module LogModule, fields ...zap.Field) *zap.Logger {
	logger := CronLoggerWithModule(ctx, module)
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	return logger
}
