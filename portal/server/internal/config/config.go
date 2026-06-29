package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var configPath = "config.yaml"

// SerConfigPath 设置配置文件路径
func SerConfigPath(configFilePath string) {
	configPath = configFilePath
}

type Config struct {
	ServerPort             int             `yaml:"ServerPort"`
	Debug                  bool            `yaml:"Debug"`                    // 调试模式
	Mysql                  MysqlConfig     `yaml:"Mysql"`
	Redis                  RedisConfig     `yaml:"Redis"`
	Logger                 LoggerConfig    `yaml:"Logger"`
	CollectBodyContentType []string        `yaml:"CollectBodyContentType"`
	GlobalRateLimit        RateLimitConfig `yaml:"GlobalRateLimit"`
	LoginJwt               LoginJwt        `yaml:"LoginJwt"`
	DataEncryptKey         string          `yaml:"DataEncryptKey"` // 数据脱敏aes加密key
	TencentSms             TencentSms      `yaml:"TencentSms"`     // 腾讯云短信配置
	TencentCloud           TencentCloud    `yaml:"TencentCloud"`   // 腾讯云配置
}

// MysqlConfig mysql配置
type MysqlConfig struct {
	Dsn             string        `yaml:"Dsn"`
	MaxIdleConns    int           `yaml:"MaxIdleConns"` // 最大空闲池
	MaxOpenConns    int           `yaml:"MaxOpenConns"` // 最大连接池
	ConnMaxLifetime time.Duration `yaml:"ConnMaxLifetime"`
	AutoMigrate     bool          `yaml:"AutoMigrate"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Mode         string   `yaml:"Mode"` // single、cluster、sentinel
	Addr         string   `yaml:"Addr"`
	Addrs        []string `yaml:"Addrs"`
	DB           int      `yaml:"DB"`
	PoolSize     int      `yaml:"PoolSize"`
	MinIdleConns int      `yaml:"MinIdleConns"`
	Password     string   `yaml:"Password"`
	MasterName   string   `yaml:"MasterName"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	LogDir       string `yaml:"LogDir"`       // 日志输出目录
	Level        string `yaml:"Level"`        // 日志输出级别
	EnableRemote bool   `yaml:"EnableRemote"` // 是否开启推送日志服务
	RemoteURL    string `yaml:"RemoteURL"`    // 日志服务地址
}

type RateLimitConfig struct {
	Rate   int           `yaml:"Rate"`   // 在 Period 单位内，允许通过的请求数
	Burst  int           `yaml:"Burst"`  // 最大桶
	Period time.Duration `yaml:"Period"` // 统计周期
}

type LoginJwt struct {
	Key     string        `yaml:"Key"`     // jwt解密key
	Expires time.Duration `yaml:"Expires"` // jwt过期数据
	Issuer  string        `yaml:"Issuer"`  // 签发人
}

type TencentSms struct {
	TemplateId              string `yaml:"TemplateId"`
	SignName                string `yaml:"SignName"`
	SdkPostDeactivateUserId string `yaml:"SdkPostDeactivateUserId"`
	IsOpen                  bool   `yaml:"IsOpen"`
}

type TencentCloud struct {
	SecretId  string `yaml:"SecretId"`
	SecretKey string `yaml:"SecretKey"`
}

// ApplyDefaults 应用默认值
func (c *Config) ApplyDefaults() {
	// ServerPort 默认值
	if c.ServerPort <= 0 {
		c.ServerPort = 8088
	}

	// Debug 默认为 false
	// Debug 模式已经通过 yaml 标签的零值为 false 处理

	// Mysql 默认值
	if c.Mysql.MaxIdleConns <= 0 {
		c.Mysql.MaxIdleConns = 10
	}
	if c.Mysql.MaxOpenConns <= 0 {
		c.Mysql.MaxOpenConns = 100
	}
	if c.Mysql.ConnMaxLifetime <= 0 {
		c.Mysql.ConnMaxLifetime = time.Hour
	}

	// Redis 默认值
	if c.Redis.PoolSize <= 0 {
		c.Redis.PoolSize = 100
	}
	if c.Redis.MinIdleConns <= 0 {
		c.Redis.MinIdleConns = 10
	}
	if c.Redis.DB < 0 {
		c.Redis.DB = 0
	}
	if c.Redis.Mode == "" {
		c.Redis.Mode = "single"
	}

	// Logger 默认值
	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}
	if c.Logger.LogDir == "" {
		c.Logger.LogDir = "logs"
	}

	// GlobalRateLimit 默认值
	if c.GlobalRateLimit.Rate <= 0 {
		c.GlobalRateLimit.Rate = 100
	}
	if c.GlobalRateLimit.Burst <= 0 {
		c.GlobalRateLimit.Burst = 200
	}
	if c.GlobalRateLimit.Period <= 0 {
		c.GlobalRateLimit.Period = time.Minute
	}

	// LoginJwt 默认值
	if c.LoginJwt.Expires <= 0 {
		c.LoginJwt.Expires = 24 * time.Hour
	}
	if c.LoginJwt.Issuer == "" {
		c.LoginJwt.Issuer = "minghe-portal"
	}

	// 确保数据加密密钥长度
	if len(c.DataEncryptKey) < 16 {
		c.DataEncryptKey = "minghe-default-encryption-key-32"
	}
}

// Validate 验证配置合法性
func (c *Config) Validate() error {
	// 验证端口范围
	if c.ServerPort < 1024 || c.ServerPort > 65535 {
		return fmt.Errorf("服务器端口必须在 1024-65535 范围内，当前值: %d", c.ServerPort)
	}

	// 验证 MySQL DSN
	if err := c.validateMysqlConfig(); err != nil {
		return fmt.Errorf("MySQL 配置验证失败: %w", err)
	}

	// 验证 Redis 配置
	if err := c.validateRedisConfig(); err != nil {
		return fmt.Errorf("Redis 配置验证失败: %w", err)
	}

	// 验证 JWT 配置
	if err := c.validateJwtConfig(); err != nil {
		return fmt.Errorf("JWT 配置验证失败: %w", err)
	}

	// 验证日志配置
	if err := c.validateLoggerConfig(); err != nil {
		return fmt.Errorf("日志配置验证失败: %w", err)
	}

	// 验证限流配置
	if err := c.validateRateLimitConfig(); err != nil {
		return fmt.Errorf("限流配置验证失败: %w", err)
	}

	// 验证腾讯云配置
	if c.TencentSms.IsOpen {
		if c.TencentSms.TemplateId == "" {
			return fmt.Errorf("开启短信服务时，TemplateId 不能为空")
		}
		if c.TencentSms.SignName == "" {
			return fmt.Errorf("开启短信服务时，SignName 不能为空")
		}
		if c.TencentCloud.SecretId == "" || c.TencentCloud.SecretKey == "" {
			return fmt.Errorf("开启短信服务时，腾讯云 SecretId 和 SecretKey 不能为空")
		}
	}

	return nil
}

// validateMysqlConfig 验证 MySQL 配置
func (c *Config) validateMysqlConfig() error {
	if c.Mysql.Dsn == "" {
		return fmt.Errorf("MySQL DSN 不能为空")
	}

	// 基本验证 DSN 格式
	if !strings.Contains(c.Mysql.Dsn, "@") {
		return fmt.Errorf("MySQL DSN 格式不正确，缺少 @ 符号")
	}
	if !strings.Contains(c.Mysql.Dsn, "(") || !strings.Contains(c.Mysql.Dsn, ")") {
		return fmt.Errorf("MySQL DSN 格式不正确，缺少主机地址括号")
	}

	// 验证连接池配置
	if c.Mysql.MaxIdleConns < 0 {
		return fmt.Errorf("MySQL 最大空闲连接数不能为负数")
	}
	if c.Mysql.MaxOpenConns < 0 {
		return fmt.Errorf("MySQL 最大连接数不能为负数")
	}
	if c.Mysql.MaxIdleConns > c.Mysql.MaxOpenConns {
		return fmt.Errorf("MySQL 最大空闲连接数不能大于最大连接数")
	}
	if c.Mysql.ConnMaxLifetime < time.Second {
		return fmt.Errorf("MySQL 连接最大生命周期至少为 1 秒")
	}

	return nil
}

// validateRedisConfig 验证 Redis 配置
func (c *Config) validateRedisConfig() error {
	// 验证 Redis 模式
	validModes := []string{"single", "cluster", "sentinel"}
	isValidMode := false
	for _, mode := range validModes {
		if c.Redis.Mode == mode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		return fmt.Errorf("Redis 模式必须是 single、cluster 或 sentinel 之一")
	}

	// 根据 Redis 模式验证配置
	switch c.Redis.Mode {
	case "single":
		if c.Redis.Addr == "" {
			return fmt.Errorf("single 模式下，Redis 地址不能为空")
		}
		// 验证地址格式
		if !strings.Contains(c.Redis.Addr, ":") {
			return fmt.Errorf("Redis 地址格式不正确，应包含端口号")
		}
	case "cluster":
		if len(c.Redis.Addrs) == 0 {
			return fmt.Errorf("cluster 模式下，Redis 节点地址列表不能为空")
		}
	case "sentinel":
		if len(c.Redis.Addrs) == 0 {
			return fmt.Errorf("sentinel 模式下，Sentinel 节点地址列表不能为空")
		}
		if c.Redis.MasterName == "" {
			return fmt.Errorf("sentinel 模式下，MasterName 不能为空")
		}
	}

	// 验证连接池配置
	if c.Redis.PoolSize < 0 {
		return fmt.Errorf("Redis 连接池大小不能为负数")
	}
	if c.Redis.MinIdleConns < 0 {
		return fmt.Errorf("Redis 最小空闲连接数不能为负数")
	}
	if c.Redis.MinIdleConns > c.Redis.PoolSize {
		return fmt.Errorf("Redis 最小空闲连接数不能大于连接池大小")
	}

	return nil
}

// validateJwtConfig 验证 JWT 配置
func (c *Config) validateJwtConfig() error {
	if c.LoginJwt.Key == "" {
		return fmt.Errorf("JWT Key 不能为空")
	}
	if len(c.LoginJwt.Key) < 16 {
		return fmt.Errorf("JWT Key 长度至少为 16 个字符")
	}
	if c.LoginJwt.Expires < time.Minute {
		return fmt.Errorf("JWT 过期时间至少为 1 分钟")
	}
	if c.LoginJwt.Issuer == "" {
		return fmt.Errorf("JWT 签发人不能为空")
	}

	return nil
}

// validateLoggerConfig 验证日志配置
func (c *Config) validateLoggerConfig() error {
	// 验证日志级别
	validLevels := map[string]bool{
		"debug":  true,
		"info":   true,
		"warn":   true,
		"error":  true,
		"panic":  true,
		"fatal":  true,
	}
	if !validLevels[strings.ToLower(c.Logger.Level)] {
		return fmt.Errorf("日志级别必须是 debug、info、warn、error、panic 或 fatal 之一")
	}

	// 验证日志目录
	if c.Logger.LogDir != "" {
		// 转换为绝对路径
		absPath, err := filepath.Abs(c.Logger.LogDir)
		if err != nil {
			return fmt.Errorf("无法获取日志目录的绝对路径: %w", err)
		}

		// 尝试创建目录
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("无法创建日志目录 %s: %w", absPath, err)
		}

		// 检查目录是否可写
		testFile := filepath.Join(absPath, ".write_test")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			return fmt.Errorf("日志目录 %s 不可写: %w", absPath, err)
		}
		_ = os.Remove(testFile)
	}

	// 验证远程日志配置
	if c.Logger.EnableRemote {
		if c.Logger.RemoteURL == "" {
			return fmt.Errorf("启用远程日志时，RemoteURL 不能为空")
		}
		// 验证 URL 格式
		if _, err := url.Parse(c.Logger.RemoteURL); err != nil {
			return fmt.Errorf("远程日志 URL 格式不正确: %w", err)
		}
	}

	return nil
}

// validateRateLimitConfig 验证限流配置
func (c *Config) validateRateLimitConfig() error {
	if c.GlobalRateLimit.Rate <= 0 {
		return fmt.Errorf("全局限流速率必须大于 0")
	}
	if c.GlobalRateLimit.Burst <= 0 {
		return fmt.Errorf("全局限流突发值必须大于 0")
	}
	if c.GlobalRateLimit.Burst < c.GlobalRateLimit.Rate {
		return fmt.Errorf("全局限流突发值不能小于速率")
	}
	if c.GlobalRateLimit.Period <= 0 {
		return fmt.Errorf("全局限流周期必须大于 0")
	}
	if c.GlobalRateLimit.Period > time.Hour {
		return fmt.Errorf("全局限流周期不能超过 1 小时")
	}

	return nil
}

func LoadConfig() (Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("配置文件加载失败: %w", err)
	}

	// 使用环境变量替换 ${VAR} 占位符
	expanded := os.ExpandEnv(string(data))

	// 继续配置文件
	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return Config{}, err
	}

	// 应用默认值
	cfg.ApplyDefaults()

	return cfg, nil
}

func LoadMysqlConfig(config Config) MysqlConfig {
	return config.Mysql
}

func LoadRedisConfig(config Config) RedisConfig {
	return config.Redis
}
