package config

import (
	"os"
	"testing"
	"time"

	"portal/internal/config"
)

// TestApplyDefaults 测试默认值应用
func TestApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   config.Config
		expected config.Config
	}{
		{
			name:   "空配置应用默认值",
			config: config.Config{},
			expected: config.Config{
				ServerPort: 8088,
				Mysql: config.MysqlConfig{
					MaxIdleConns:    10,
					MaxOpenConns:    100,
					ConnMaxLifetime: time.Hour,
				},
				Redis: config.RedisConfig{
					Mode:         "single",
					DB:           0,
					PoolSize:     100,
					MinIdleConns: 10,
				},
				Logger: config.LoggerConfig{
					Level:  "info",
					LogDir: "logs",
				},
				GlobalRateLimit: config.RateLimitConfig{
					Rate:   100,
					Burst:  200,
					Period: time.Minute,
				},
				LoginJwt: config.LoginJwt{
					Expires: 24 * time.Hour,
					Issuer:  "minghe-portal",
				},
				DataEncryptKey: "minghe-default-encryption-key-32",
			},
		},
		{
			name: "部分配置值不应用默认值",
			config: config.Config{
				ServerPort: 9000,
				Mysql: config.MysqlConfig{
					Dsn:             "test",
					MaxIdleConns:    20,
					MaxOpenConns:    200,
					ConnMaxLifetime: 2 * time.Hour,
				},
			},
			expected: config.Config{
				ServerPort: 9000,
				Mysql: config.MysqlConfig{
					Dsn:             "test",
					MaxIdleConns:    20,
					MaxOpenConns:    200,
					ConnMaxLifetime: 2 * time.Hour,
				},
				Redis: config.RedisConfig{
					Mode:         "single",
					DB:           0,
					PoolSize:     100,
					MinIdleConns: 10,
				},
				Logger: config.LoggerConfig{
					Level:  "info",
					LogDir: "logs",
				},
				GlobalRateLimit: config.RateLimitConfig{
					Rate:   100,
					Burst:  200,
					Period: time.Minute,
				},
				LoginJwt: config.LoginJwt{
					Expires: 24 * time.Hour,
					Issuer:  "minghe-portal",
				},
				DataEncryptKey: "minghe-default-encryption-key-32",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.ApplyDefaults()

			if tt.config.ServerPort != tt.expected.ServerPort {
				t.Errorf("ServerPort = %v, want %v", tt.config.ServerPort, tt.expected.ServerPort)
			}
			if tt.config.Mysql.MaxIdleConns != tt.expected.Mysql.MaxIdleConns {
				t.Errorf("Mysql.MaxIdleConns = %v, want %v", tt.config.Mysql.MaxIdleConns, tt.expected.Mysql.MaxIdleConns)
			}
			if tt.config.Redis.Mode != tt.expected.Redis.Mode {
				t.Errorf("Redis.Mode = %v, want %v", tt.config.Redis.Mode, tt.expected.Redis.Mode)
			}
			if tt.config.Logger.Level != tt.expected.Logger.Level {
				t.Errorf("Logger.Level = %v, want %v", tt.config.Logger.Level, tt.expected.Logger.Level)
			}
		})
	}
}

// TestValidatePortRange 测试端口范围验证
func TestValidatePortRange(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"有效端口8088", 8088, false},
		{"有效端口1024", 1024, false},
		{"有效端口65535", 65535, false},
		{"无效端口1023", 1023, true},
		{"无效端口65536", 65536, true},
		{"无效端口0", 0, true},
		{"无效端口-1", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
				cfg := createValidConfig(t)
			cfg.ServerPort = tt.port

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateMysqlDSN 测试 MySQL DSN 验证
func TestValidateMysqlDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{"有效DSN", "user:pass@tcp(localhost:3306)/db", false},
		{"有效DSN带参数", "user:pass@tcp(127.0.0.1:3306)/db?charset=utf8mb4", false},
		{"缺少@", "user:passtcp(localhost:3306)/db", true},
		{"缺少主机地址", "user:pass@/db", true},
		{"空DSN", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createValidConfig(t)
			cfg.Mysql.Dsn = tt.dsn

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateJwtKey 测试 JWT Key 验证
func TestValidateJwtKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"有效Key", "this-is-a-very-long-jwt-key-for-testing", false},
		{"最小长度Key", "1234567890123456", false},
		{"空Key", "", true},
		{"过短Key", "short", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createValidConfig(t)
			cfg.LoginJwt.Key = tt.key

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateLogLevel 测试日志级别验证
func TestValidateLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"debug级别", "debug", false},
		{"info级别", "info", false},
		{"warn级别", "warn", false},
		{"error级别", "error", false},
		{"无效级别", "invalid", true},
		{"空级别", "", false}, // 会应用默认值
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createValidConfig(t)
			if tt.level != "" {
				cfg.Logger.Level = tt.level
			}

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateRedisMode 测试 Redis 模式验证
func TestValidateRedisMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{"single模式", "single", false},
		{"cluster模式", "cluster", false},
		{"sentinel模式", "sentinel", false},
		{"无效模式", "invalid", true},
		{"空模式", "", false}, // 会应用默认值
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createValidConfig(t)
			if tt.mode != "" {
				cfg.Redis.Mode = tt.mode
			}

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// createValidConfig 创建一个有效的配置用于测试
func createValidConfig(t *testing.T) config.Config {
	tempDir := t.TempDir()

	return config.Config{
		ServerPort: 8088,
		Mysql: config.MysqlConfig{
			Dsn:             "user:pass@tcp(localhost:3306)/db",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: time.Hour,
		},
		Redis: config.RedisConfig{
			Mode:     "single",
			Addr:     "localhost:6379",
			PoolSize: 100,
		},
		LoginJwt: config.LoginJwt{
			Key:     "this-is-a-very-long-jwt-key-for-testing",
			Expires: time.Hour,
			Issuer:  "test",
		},
		Logger: config.LoggerConfig{
			Level:  "info",
			LogDir: tempDir,
		},
		GlobalRateLimit: config.RateLimitConfig{
			Rate:   100,
			Burst:  200,
			Period: time.Minute,
		},
	}
}

// TestValidateLogDirWritable 测试日志目录可写性验证
func TestValidateLogDirWritable(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	cfg := createValidConfig(t)
	cfg.Logger.LogDir = tempDir

	// 应该通过验证
	if err := cfg.Validate(); err != nil {
		t.Errorf("Valid log directory should pass validation: %v", err)
	}

	// 测试不可写的目录
	readonlyDir := t.TempDir()
	cfg.Logger.LogDir = readonlyDir
	// 在某些系统上，我们无法真正设置目录为只读
	// 这里只测试目录存在的情况

	if err := cfg.Validate(); err != nil {
		t.Errorf("Writable log directory should pass validation: %v", err)
	}
}

// TestValidateTencentSmsConfig 测试腾讯云短信配置验证
func TestValidateTencentSmsConfig(t *testing.T) {
	tests := []struct {
		name    string
		isOpen  bool
		config  config.TencentSms
		cloud   config.TencentCloud
		wantErr bool
	}{
		{
			name:   "关闭短信服务",
			isOpen: false,
			config: config.TencentSms{},
			cloud:  config.TencentCloud{},
		},
		{
			name:   "开启短信服务-完整配置",
			isOpen: true,
			config: config.TencentSms{
				TemplateId: "123456",
				SignName:   "test",
			},
			cloud: config.TencentCloud{
				SecretId:  "secret_id",
				SecretKey: "secret_key",
			},
		},
		{
			name:   "开启短信服务-缺少模板ID",
			isOpen: true,
			config: config.TencentSms{
				SignName: "test",
			},
			cloud: config.TencentCloud{
				SecretId:  "secret_id",
				SecretKey: "secret_key",
			},
			wantErr: true,
		},
		{
			name:   "开启短信服务-缺少签名",
			isOpen: true,
			config: config.TencentSms{
				TemplateId: "123456",
			},
			cloud: config.TencentCloud{
				SecretId:  "secret_id",
				SecretKey: "secret_key",
			},
			wantErr: true,
		},
		{
			name:   "开启短信服务-缺少密钥",
			isOpen: true,
			config: config.TencentSms{
				TemplateId: "123456",
				SignName:   "test",
			},
			cloud:   config.TencentCloud{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createValidConfig(t)
			cfg.TencentSms = tt.config
			cfg.TencentSms.IsOpen = tt.isOpen
			cfg.TencentCloud = tt.cloud

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateRateLimitConfig 测试限流配置验证
func TestValidateRateLimitConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  config.RateLimitConfig
		wantErr bool
	}{
		{"有效配置", config.RateLimitConfig{Rate: 100, Burst: 200, Period: time.Minute}, false},
		{"速率为0", config.RateLimitConfig{Rate: 0, Burst: 200, Period: time.Minute}, true},
		{"突发值为0", config.RateLimitConfig{Rate: 100, Burst: 0, Period: time.Minute}, true},
		{"周期为0", config.RateLimitConfig{Rate: 100, Burst: 200, Period: 0}, true},
		{"突发值小于速率", config.RateLimitConfig{Rate: 200, Burst: 100, Period: time.Minute}, true},
		{"周期超过1小时", config.RateLimitConfig{Rate: 100, Burst: 200, Period: 2 * time.Hour}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createValidConfig(t)
			cfg.GlobalRateLimit = tt.config

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestLoadConfigWithEnvVars 测试环境变量替换
func TestLoadConfigWithEnvVars(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("TEST_SERVER_PORT", "9999")
	os.Setenv("TEST_MYSQL_DSN", "test:user:pass@tcp(localhost:3306)/db")

	// 创建临时配置文件
	configContent := `
ServerPort: ${TEST_SERVER_PORT}
Mysql:
  Dsn: ${TEST_MYSQL_DSN}
  MaxIdleConns: 10
  MaxOpenConns: 100
  ConnMaxLifetime: 1h
Redis:
  Mode: single
  Addr: localhost:6379
  PoolSize: 100
  MinIdleConns: 10
  DB: 0
Logger:
  Level: info
  LogDir: logs
  EnableRemote: false
LoginJwt:
  Key: this-is-a-very-long-jwt-key-for-testing
  Expires: 24h
  Issuer: test
GlobalRateLimit:
  Rate: 100
  Burst: 200
  Period: 1m
DataEncryptKey: minghe-default-encryption-key-32
TencentSms:
  IsOpen: false
TencentCloud:
  SecretId: ""
  SecretKey: ""
`

	tempFile, err := os.CreateTemp(t.TempDir(), "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tempFile.Close()

	// 设置配置文件路径
	config.SerConfigPath(tempFile.Name())

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证环境变量替换
	if cfg.ServerPort != 9999 {
		t.Errorf("Expected ServerPort 9999, got %d", cfg.ServerPort)
	}
	if cfg.Mysql.Dsn != "test:user:pass@tcp(localhost:3306)/db" {
		t.Errorf("Expected DSN with env var, got %s", cfg.Mysql.Dsn)
	}

	// 清理环境变量
	os.Unsetenv("TEST_SERVER_PORT")
	os.Unsetenv("TEST_MYSQL_DSN")
}
