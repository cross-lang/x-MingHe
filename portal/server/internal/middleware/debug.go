package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"portal/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// DebugMiddleware 调试模式中间件
// 在调试模式下美化 JSON 输出
func DebugMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只在调试模式下执行
		if !cfg.Debug {
			c.Next()
			return
		}

		// 创建响应写入器来捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// 只处理 JSON 响应
		contentType := c.Writer.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			// 解析 JSON 以美化输出
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, blw.body.Bytes(), "", "  "); err == nil {
				c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", prettyJSON.Len()))
				c.Writer.WriteHeader(c.Writer.Status())
				c.Writer.Write(prettyJSON.Bytes())
			}
		}
	}
}

// bodyLogWriter 用于捕获响应内容的写入器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
