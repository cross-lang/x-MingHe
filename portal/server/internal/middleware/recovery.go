package middleware

import (
	"portal/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogWriter struct {
	logger *zap.Logger
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.logger.Error(string(p))
	return len(p), nil
}

// Recovery 中间件
func Recovery() gin.HandlerFunc {
	writer := LogWriter{logger: log.Logger}
	return gin.RecoveryWithWriter(&writer)
}
