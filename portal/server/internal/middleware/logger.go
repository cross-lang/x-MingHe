package middleware

import (
	"bytes"
	"io"
	"portal/internal/pkg/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// contentType前缀匹配
func contentTypePrefixMatch(contentTypeWhitelist []string, contentType string) bool {
	if len(contentTypeWhitelist) == 0 {
		return true
	}
	if contentType == "" {
		return false
	}
	ct := strings.ToLower(contentType)
	for _, v := range contentTypeWhitelist {
		if strings.HasPrefix(ct, strings.ToLower(v)) {
			return true
		}
	}
	return false
}

// ResponseBodyWriter 用于捕获响应内容
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware(contentTypeWhitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger := log.WithContext(c.Request.Context())

		// 读取请求 body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 捕获响应 body
		writer := &ResponseBodyWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = writer

		// 请求日志
		reqFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
		}
		isCollectReqBody := contentTypePrefixMatch(contentTypeWhitelist, c.Request.Header.Get("Content-Type"))
		if isCollectReqBody {
			reqFields = append(reqFields, zap.ByteString("request_body", reqBody))
		}
		logger.Info("request", reqFields...)

		// 处理请求
		c.Next()

		duration := time.Since(start)

		// 响应日志
		respFields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		}
		isCollectRespBody := contentTypePrefixMatch(contentTypeWhitelist, c.Writer.Header().Get("Content-Type"))
		if isCollectRespBody {
			respFields = append(respFields, zap.ByteString("body", writer.body.Bytes()))
		}
		logger.Info("response", respFields...)
	}
}
