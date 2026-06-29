package middleware

import (
	"portal/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceMiddleware 链路追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 写入 request context
		ctx := log.ContextWithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 设置 gin context 中的 request_id，用于响应
		c.Set("request_id", traceID)

		// 设置响应 header
		c.Writer.Header().Set("X-Trace-ID", traceID)

		c.Next()
	}
}
