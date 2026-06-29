package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"portal/internal/config"
	"portal/internal/constant"
	"portal/internal/pkg"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/ginx"
	"portal/internal/pkg/log"
	"portal/internal/repository"
	"portal/internal/types"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// LoginAuthMiddleware 用户登录校验中间件
func LoginAuthMiddleware(userRepo *repository.UserRepo, cache *repository.Cache, conf config.LoginJwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 解析jwt
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := parseToken(tokenStr, conf.Key)
		if err != nil {
			log.WithContext(c.Request.Context()).Error("解析jwt失败", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 检查用户登录状态
		status, err := cache.GetLoginStatus(c.Request.Context(), claims.UserId, tokenStr)
		if err != nil {
			log.WithContext(c.Request.Context()).Error("检查用户登录状态失败", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !status {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 获取用户信息
		userInfo, err := userRepo.GetById(c.Request.Context(), claims.UserId)
		if err != nil {
			log.WithContext(c.Request.Context()).Error("获取用户失败", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 黑名单拦截逻辑
		if userInfo.BlockStatus == constant.UserBlockStatusYes && isWriteOperation(c) {
			log.WithContext(c.Request.Context()).Error("您暂时不能使用该功能，如有疑问请联系园区管理员", zap.Any("userInfo", userInfo))
			ginx.FailResponseWithCode(c, errorx.UserBlockedError.Code, errorx.UserBlockedError.Message)
			return
		}

		// 保存用户信息到上下文
		ginx.SetUserToCtx(c, userInfo)

		// 保存用户信息到标准上下文
		newCtx := pkg.SetUserToCtx(c.Request.Context(), userInfo)
		// 重写请求中的标准上下文
		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	}
}

// isWriteOperation 判断是否为需要拦截的写操作（创建/提交类）
func isWriteOperation(c *gin.Context) bool {
	method := c.Request.Method
	fullPath := c.FullPath()

	if fullPath == "" {
		return false
	}

	if method != http.MethodPost && method != http.MethodPut && method != http.MethodDelete {
		return false
	}

	// 直接匹配原始路由模板
	writePatterns := []string{
		"/v1/services/:service_id/orders/place",
		"/v1/services/orders/:order_id/update",
		"/v1/services/orders/:order_id/clone",
		"/v1/services/orders/:order_id/delete",
		"/v1/services/orders/:order_id/recall",
		"/v1/activities/:activity_id/sign-up/execute",
		"/v1/activities/:activity_id/sign-up/cancel",
	}

	for _, pattern := range writePatterns {
		if fullPath == pattern {
			return true
		}
	}

	return false
}

// parseToken 解析jwt
func parseToken(jwtToken, secret string) (*types.UserLoginJwtClaims, error) {
	var claims types.UserLoginJwtClaims
	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		// 验证算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}
