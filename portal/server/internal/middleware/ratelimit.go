package middleware

import (
	"portal/internal/config"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/ginx"
	"portal/internal/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"go.uber.org/zap"
)

// GlobalRateLimitMiddleware 全局限流器按ip限流
func GlobalRateLimitMiddleware(rateLimiter *redis_rate.Limiter, conf config.RateLimitConfig) gin.HandlerFunc {
	return RateLimitMiddleware(rateLimiter, conf, func(c *gin.Context) string {
		return "rl:" + c.ClientIP()
	})
}

// SpecializedRateLimitMiddleware 特定场景的限流器按接口ip+path限流
func SpecializedRateLimitMiddleware(rateLimiter *redis_rate.Limiter, period string, rate int) gin.HandlerFunc {
	duration, err := time.ParseDuration(period)
	if err != nil {
		duration = time.Second
	}
	return RateLimitMiddleware(rateLimiter, config.RateLimitConfig{
		Rate:   rate,
		Burst:  rate,
		Period: duration,
	}, func(c *gin.Context) string {
		return "rl:" + c.ClientIP() + ":" + c.FullPath()
	})
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(rateLimiter *redis_rate.Limiter, conf config.RateLimitConfig, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c) // 全局限流器按ip限流
		res, err := rateLimiter.Allow(c.Request.Context(), key, redis_rate.Limit{
			Rate:   conf.Rate,
			Burst:  conf.Burst,
			Period: conf.Period,
		})
		if err != nil {
			log.WithContext(c.Request.Context()).Error("限流器和token发生错误", zap.Error(err))
			c.Next()
			return
		}
		if res.Allowed == 0 {
			c.Header("Retry-After", res.RetryAfter.String())
			ginx.FailureResponse(c, errorx.RateLimitError)
			return
		}
		c.Next()
	}
}
