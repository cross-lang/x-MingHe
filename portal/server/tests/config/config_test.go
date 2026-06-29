package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"portal/internal/config"
)

func TestSerConfigPath(t *testing.T) {
	t.Run("sets config path", func(t *testing.T) {
		originalPath := "test_config.yaml"
		config.SerConfigPath(originalPath)

		// Cannot directly test the internal configPath variable
		// Just verify the function doesn't panic
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("load valid config", func(t *testing.T) {
		// Create a temporary config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "test_config.yaml")
		configContent := `
ServerPort: 8088
Mysql:
  Dsn: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
  MaxIdleConns: 10
  MaxOpenConns: 100
  ConnMaxLifetime: 1h
  AutoMigrate: false
Redis:
  Mode: "single"
  Addr: "127.0.0.1:6379"
  DB: 0
  PoolSize: 10
  MinIdleConns: 5
  Password: ""
  MasterName: ""
Logger:
  LogDir: "log"
  Level: "info"
  EnableRemote: false
  RemoteURL: ""
GlobalRateLimit:
  Rate: 100
  Burst: 200
  Period: 1m
LoginJwt:
  Key: "testkey1234567"
  Expires: 2160h
  Issuer: "test"
DataEncryptKey: "test-encrypt-key-123456789"
TencentSms:
  TemplateId: "123456"
  SignName: "TestSign"
  SdkPostDeactivateUserId: ""
  IsOpen: false
TencentCloud:
  SecretId: "test_secret_id"
  SecretKey: "test_secret_key"
CollectBodyContentType:
  - "application/json"
  - "application/x-www-form-urlencoded"
`
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		// Test loading the config
		config.SerConfigPath(configPath)
		cfg, err := config.LoadConfig()
		if err != nil {
			t.Errorf("LoadConfig() failed: %v", err)
			return
		}

		// Verify loaded values
		if cfg.ServerPort != 8088 {
			t.Errorf("ServerPort = %d, want 8088", cfg.ServerPort)
		}
		if cfg.Mysql.Dsn != "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" {
			t.Errorf("Mysql.Dsn incorrect")
		}
		if cfg.Redis.Addr != "127.0.0.1:6379" {
			t.Errorf("Redis.Addr = %s, want 127.0.0.1:6379", cfg.Redis.Addr)
		}
		if cfg.Logger.Level != "info" {
			t.Errorf("Logger.Level = %s, want info", cfg.Logger.Level)
		}
	})
}

func TestLoadConfigWithEnvVars(t *testing.T) {
	t.Run("load config with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("SERVER_PORT", "9090")
		os.Setenv("MYSQL_DSN", "root:env@tcp(localhost:3306)/envdb?charset=utf8mb4")
		os.Setenv("REDIS_ADDR", "localhost:6380")
		defer os.Unsetenv("SERVER_PORT")
		defer os.Unsetenv("MYSQL_DSN")
		defer os.Unsetenv("REDIS_ADDR")

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "test_env_config.yaml")
		configContent := `
ServerPort: ${SERVER_PORT}
Mysql:
  Dsn: ${MYSQL_DSN}
  MaxIdleConns: 10
  MaxOpenConns: 100
  ConnMaxLifetime: 1h
  AutoMigrate: false
Redis:
  Mode: "single"
  Addr: ${REDIS_ADDR}
  DB: 0
  PoolSize: 10
  MinIdleConns: 5
  Password: ""
  MasterName: ""
Logger:
  LogDir: "log"
  Level: "debug"
  EnableRemote: false
  RemoteURL: ""
GlobalRateLimit:
  Rate: 100
  Burst: 200
  Period: 1m
LoginJwt:
  Key: "testkey1234567"
  Expires: 2160h
  Issuer: "test"
DataEncryptKey: "test-encrypt-key-123456789"
TencentSms:
  TemplateId: "123456"
  SignName: "TestSign"
  SdkPostDeactivateUserId: ""
  IsOpen: false
TencentCloud:
  SecretId: "test_secret_id"
  SecretKey: "test_secret_key"
CollectBodyContentType: []
`
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		config.SerConfigPath(configPath)
		cfg, err := config.LoadConfig()
		if err != nil {
			t.Errorf("LoadConfig() with env vars failed: %v", err)
			return
		}

		// Verify env vars were expanded
		if cfg.ServerPort != 9090 {
			t.Errorf("ServerPort with env var = %d, want 9090", cfg.ServerPort)
		}
		if cfg.Mysql.Dsn != "root:env@tcp(localhost:3306)/envdb?charset=utf8mb4" {
			t.Errorf("Mysql.Dsn with env var not expanded correctly")
		}
		if cfg.Redis.Addr != "localhost:6380" {
			t.Errorf("Redis.Addr with env var = %s, want localhost:6380", cfg.Redis.Addr)
		}
	})
}

func TestLoadConfigMissingFile(t *testing.T) {
	t.Run("load non-existent config file", func(t *testing.T) {
		config.SerConfigPath("non_existent_config.yaml")
		_, err := config.LoadConfig()

		if err == nil {
			t.Error("LoadConfig() should fail for non-existent file")
		}
	})
}

func TestLoadConfigInvalidYaml(t *testing.T) {
	t.Run("load invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "invalid_config.yaml")
		invalidYaml := `
ServerPort: 8088
Mysql:
  Dsn: "test"
  MaxIdleConns: this_is_not_a_number
`
		err := os.WriteFile(configPath, []byte(invalidYaml), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		config.SerConfigPath(configPath)
		_, err = config.LoadConfig()

		if err == nil {
			t.Error("LoadConfig() should fail for invalid YAML")
		}
	})
}

func TestLoadMysqlConfig(t *testing.T) {
	t.Run("extract Mysql config", func(t *testing.T) {
		cfg := config.Config{
			Mysql: config.MysqlConfig{
				Dsn:             "test:dns",
				MaxIdleConns:    5,
				MaxOpenConns:    50,
				ConnMaxLifetime: 30 * 60 * 1000000000,
				AutoMigrate:     true,
			},
		}

		mysqlCfg := config.LoadMysqlConfig(cfg)

		if mysqlCfg.Dsn != "test:dns" {
			t.Errorf("LoadMysqlConfig() Dsn = %s, want test:dns", mysqlCfg.Dsn)
		}
		if mysqlCfg.MaxIdleConns != 5 {
			t.Errorf("LoadMysqlConfig() MaxIdleConns = %d, want 5", mysqlCfg.MaxIdleConns)
		}
	})
}

func TestLoadRedisConfig(t *testing.T) {
	t.Run("extract Redis config", func(t *testing.T) {
		cfg := config.Config{
			Redis: config.RedisConfig{
				Mode:         "cluster",
				Addr:         "localhost:6379",
				Addrs:        []string{"addr1:6379", "addr2:6379"},
				DB:           1,
				PoolSize:     20,
				MinIdleConns: 10,
				Password:     "password",
				MasterName:   "mymaster",
			},
		}

		redisCfg := config.LoadRedisConfig(cfg)

		if redisCfg.Mode != "cluster" {
			t.Errorf("LoadRedisConfig() Mode = %s, want cluster", redisCfg.Mode)
		}
		if redisCfg.Addr != "localhost:6379" {
			t.Errorf("LoadRedisConfig() Addr = %s, want localhost:6379", redisCfg.Addr)
		}
		if redisCfg.DB != 1 {
			t.Errorf("LoadRedisConfig() DB = %d, want 1", redisCfg.DB)
		}
	})
}

func TestConfigDefaults(t *testing.T) {
	t.Run("minimal config with defaults", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "minimal_config.yaml")
		minimalYaml := `
ServerPort: 8088
Mysql:
  Dsn: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
  MaxIdleConns: 10
  MaxOpenConns: 100
  ConnMaxLifetime: 1h
  AutoMigrate: false
Redis:
  Mode: "single"
  Addr: "127.0.0.1:6379"
  DB: 0
  PoolSize: 10
  MinIdleConns: 5
  Password: ""
  MasterName: ""
Logger:
  LogDir: "log"
  Level: "info"
  EnableRemote: false
  RemoteURL: ""
GlobalRateLimit:
  Rate: 100
  Burst: 200
  Period: 1m
LoginJwt:
  Key: "testkey"
  Expires: 2160h
  Issuer: "minghe"
DataEncryptKey: "test-key"
TencentSms:
  TemplateId: "123456"
  SignName: "TestSign"
  SdkPostDeactivateUserId: ""
  IsOpen: false
TencentCloud:
  SecretId: ""
  SecretKey: ""
CollectBodyContentType: []
`
		err := os.WriteFile(configPath, []byte(minimalYaml), 0644)
		if err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		config.SerConfigPath(configPath)
		cfg, err := config.LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() failed: %v", err)
		}

		// Verify all required fields are loaded
		if cfg.ServerPort == 0 {
			t.Error("ServerPort should not be zero")
		}
		if cfg.Mysql.Dsn == "" {
			t.Error("Mysql.Dsn should not be empty")
		}
		if cfg.Redis.Addr == "" {
			t.Error("Redis.Addr should not be empty")
		}
		if cfg.Logger.LogDir == "" {
			t.Error("Logger.LogDir should not be empty")
		}
	})
}
