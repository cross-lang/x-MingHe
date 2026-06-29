package repository

import (
	"context"
	"fmt"
	"portal/internal/pkg/errorx"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 缓存访问层
type Cache struct {
	RedisDB redis.UniversalClient
}

// AcquireTaskLock 获取定时任务分布式锁
// taskKey: 任务唯一标识
// ttl: 预估任务执行时间（锁过期时间）
func (c *Cache) AcquireTaskLock(ctx context.Context, taskKey string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("task:lock:%s", taskKey)
	locked, err := c.RedisDB.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, err
	}
	return locked, nil
}

// ReleaseTaskLock 释放定时任务分布式锁
func (c *Cache) ReleaseTaskLock(ctx context.Context, taskKey string) error {
	key := fmt.Sprintf("task:lock:%s", taskKey)
	return c.RedisDB.Del(ctx, key).Err()
}

// SetRegisterCode 设置用户注册验证码缓存
func (c *Cache) SetRegisterCode(ctx context.Context, phoneNumber, code string) error {
	key := fmt.Sprintf("user:register:code:%s", phoneNumber)
	effectiveKey := fmt.Sprintf("user:phone:effective:%s", phoneNumber)

	keyExists, _ := c.RedisDB.Exists(ctx, effectiveKey).Result()
	if keyExists == 1 {
		times, err := c.RedisDB.Incr(ctx, effectiveKey).Result()
		if err != nil {
			return err
		}
		if times >= 5 {
			return errorx.SmsLimitError
		}
	} else {
		if err := c.RedisDB.Set(ctx, effectiveKey, 1, 10*time.Minute).Err(); err != nil {
			return err
		}
	}
	// 有效期5分钟
	return c.RedisDB.Set(ctx, key, code, 5*time.Minute).Err()
}

// GetRegisterCode 获取用户注册验证码缓存
func (c *Cache) GetRegisterCode(ctx context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf("user:register:code:%s", phoneNumber)
	code, err := c.RedisDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

// ClearRegisterCode 清除用户注册验证码缓存
func (c *Cache) ClearRegisterCode(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("user:register:code:%s", phoneNumber)
	return c.RedisDB.Del(ctx, key).Err()
}

// SetLoginStatus 设置用户登录状态
func (c *Cache) SetLoginStatus(ctx context.Context, userId uint32, ttl time.Duration, token string) error {
	key := fmt.Sprintf("user:login:status:%d", userId)
	return c.RedisDB.Set(ctx, key, token, ttl).Err()
}

// GetLoginStatus 获取用户登录状态
func (c *Cache) GetLoginStatus(ctx context.Context, userId uint32, token string) (bool, error) {
	key := fmt.Sprintf("user:login:status:%d", userId)
	storedToken, err := c.RedisDB.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return token == storedToken, nil
}

// ClearLoginStatus 移除用户登录状态
func (c *Cache) ClearLoginStatus(ctx context.Context, userId uint32) error {
	key := fmt.Sprintf("user:login:status:%d", userId)
	return c.RedisDB.Del(ctx, key).Err()
}

// SetLoginCode 设置用户登录验证码缓存
func (c *Cache) SetLoginCode(ctx context.Context, phoneNumber, code string) error {
	key := fmt.Sprintf("user:login:code:%s", phoneNumber)
	effectiveKey := fmt.Sprintf("user:phone:effective:%s", phoneNumber)

	keyExists, _ := c.RedisDB.Exists(ctx, effectiveKey).Result()
	if keyExists == 1 {
		times, err := c.RedisDB.Incr(ctx, effectiveKey).Result()
		if err != nil {
			return err
		}
		if times >= 5 {
			return errorx.SmsLimitError
		}
	} else {
		if err := c.RedisDB.Set(ctx, effectiveKey, 1, 10*time.Minute).Err(); err != nil {
			return err
		}
	}
	// 有效期5分钟
	return c.RedisDB.Set(ctx, key, code, 5*time.Minute).Err()
}

// GetLoginCode 获取用户登录验证码缓存
func (c *Cache) GetLoginCode(ctx context.Context, phoneNumber string) (string, error) {
	key := fmt.Sprintf("user:login:code:%s", phoneNumber)
	code, err := c.RedisDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

// GetPhoneTimes 获取手机号验证码发送次数
func (c *Cache) GetPhoneTimes(ctx context.Context, phoneNumber string) (int, error) {
	effectiveKey := fmt.Sprintf("user:phone:effective:%s", phoneNumber)
	times, err := c.RedisDB.Get(ctx, effectiveKey).Int()
	if err != nil {
		return 0, err
	}
	return times, nil
}

// ClearLoginCode 清除用户登录验证码缓存
func (c *Cache) ClearLoginCode(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("user:login:code:%s", phoneNumber)
	return c.RedisDB.Del(ctx, key).Err()
}
