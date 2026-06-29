package log

import (
	"reflect"
	"strings"
	"testing"
	"portal/internal/pkg/log"
)

// TestMaskPhone 测试手机号脱敏
func TestMaskPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{"标准手机号", "13812345678", "138****5678"},
		{"带空格手机号", "138 1234 5678", "138****5678"},
		{"带连字符手机号", "138-1234-5678", "138****5678"},
		{"过短", "123456", "****"},
		{"空字符串", "", ""},
		{"非数字", "abcdefghijk", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.phone, log.MaskTypePhone)
			if result != tt.expected {
				t.Errorf("MaskPhone() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskEmail 测试邮箱脱敏
func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{"标准邮箱", "example@test.com", "e****@test.com"},
		{"长用户名", "username@example.com", "u****@example.com"},
		{"单字符用户名", "a@test.com", "****@test.com"},
		{"带数字", "user123@test.com", "u****@test.com"},
		{"空字符串", "", "****@****.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.email, log.MaskTypeEmail)
			if result != tt.expected {
				t.Errorf("MaskEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskIDCard 测试身份证号脱敏
func TestMaskIDCard(t *testing.T) {
	tests := []struct {
		name     string
		idCard   string
		expected string
	}{
		{"18位身份证", "110101199001011234", "110101********1234"},
		{"18位身份证带X", "11010119900101123X", "110101********123X"},
		{"18位身份证带空格", "110101 1990 0101 1234", "110101********1234"},
		{"15位身份证", "110101900101123", "110101*******123"},
		{"过短", "12345678", "****"},
		{"空字符串", "", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.idCard, log.MaskTypeIDCard)
			if result != tt.expected {
				t.Errorf("MaskIDCard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskBankCard 测试银行卡号脱敏
func TestMaskBankCard(t *testing.T) {
	tests := []struct {
		name     string
		cardNum  string
		expected string
	}{
		{"16位银行卡", "6225880212345678", "622588****5678"},
		{"19位银行卡", "6225880212345678901", "622588****678901"},
		{"带空格", "6225 8802 1234 5678", "622588****5678"},
		{"过短", "123456789", "****"},
		{"空字符串", "", "****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.cardNum, log.MaskTypeBankCard)
			if result != tt.expected {
				t.Errorf("MaskBankCard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskPassword 测试密码脱敏
func TestMaskPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{"普通密码", "password123", "******"},
		{"空密码", "", "******"},
		{"短密码", "123", "******"},
		{"长密码", "this_is_a_very_long_password", "******"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.password, log.MaskTypePassword)
			if result != tt.expected {
				t.Errorf("MaskPassword() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskName 测试姓名脱敏
func TestMaskName(t *testing.T) {
	tests := []struct {
		name     string
		realName string
		expected string
	}{
		{"单字姓名", "张", "*"},
		{"双字姓名", "张三", "张*"},
		{"三字姓名", "张三丰", "张**"},
		{"四字姓名", "欧阳娜娜", "欧阳**"},
		{"复姓双字", "欧阳修", "欧阳*"},
		{"复姓三字", "司马相如", "司马**"},
		{"空字符串", "", "*"},
		{"英文姓名", "John", "J***"},
		{"英文姓名双字", "John Smith", "J*** *****"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.realName, log.MaskTypeName)
			if result != tt.expected {
				t.Errorf("MaskName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskAddress 测试地址脱敏
func TestMaskAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{"标准地址", "北京市朝阳区某某路123号", false},
		{"长地址", "上海市浦东新区张江高科技园区科苑路88号", false},
		{"短地址", "北京市", true},
		{"空字符串", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.Mask(tt.address, log.MaskTypeAddress)
			if tt.wantErr {
				if result != "****" {
					t.Errorf("MaskAddress() = %v, want ****", result)
				}
			} else {
				if len(result) < 6 || !strings.Contains(result, "...") {
					t.Errorf("MaskAddress() = %v, want format like '...'", result)
				}
			}
		})
	}
}

// TestMaskField 测试字段名自动脱敏
func TestMaskField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     string
		expected  string
	}{
		{"手机号字段", "phone_number", "13812345678", "138****5678"},
		{"手机号字段2", "mobile", "15987654321", "159****4321"},
		{"邮箱字段", "email", "user@example.com", "u****@example.com"},
		{"密码字段", "password", "mypassword123", "******"},
		{"密码字段2", "pwd", "mypassword123", "******"},
		{"身份证字段", "id_card", "110101199001011234", "110101********1234"},
		{"银行卡字段", "bank_card", "6225880212345678", "622588****5678"},
		{"姓名字段", "user_name", "张三", "张*"},
		{"地址字段", "address", "北京市朝阳区", "北京...阳区"},
		{"普通字段", "description", "这是一个描述", "这是一个描述"},
		{"空值", "phone", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.MaskField(tt.fieldName, tt.value)
			if result != tt.expected {
				t.Errorf("MaskField() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestMaskMap 测试Map脱敏
func TestMaskMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "基本脱敏",
			input: map[string]interface{}{
				"username": "admin",
				"phone":    "13812345678",
				"email":    "admin@example.com",
				"password": "admin123",
			},
			expected: map[string]interface{}{
				"username": "admin",
				"phone":    "138****5678",
				"email":    "a****@example.com",
				"password": "******",
			},
		},
		{
			name: "嵌套Map",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "张三",
					"phone":    "15987654321",
					"id_card":  "110101199001011234",
				},
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "张*",
					"phone":    "159****4321",
					"id_card":  "110101********1234",
				},
			},
		},
		{
			name:     "空Map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.MaskMap(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MaskMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestLogLevel 测试日志级别
func TestLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   log.LogLevel
		valid   bool
		string  string
	}{
		{"Debug级别", log.LogLevelDebug, true, "debug"},
		{"Info级别", log.LogLevelInfo, true, "info"},
		{"Warning级别", log.LogLevelWarning, true, "warn"},
		{"Error级别", log.LogLevelError, true, "error"},
		{"Critical级别", log.LogLevelCritical, true, "critical"},
		{"无效级别", log.LogLevel("invalid"), false, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.level.IsValid() != tt.valid {
				t.Errorf("LogLevel.IsValid() = %v, want %v", tt.level.IsValid(), tt.valid)
			}
			if tt.level.String() != tt.string {
				t.Errorf("LogLevel.String() = %v, want %v", tt.level.String(), tt.string)
			}
		})
	}
}

// TestParseLevel 测试日志级别解析
func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected log.LogLevel
	}{
		{"debug", "debug", log.LogLevelDebug},
		{"DEBUG", "DEBUG", log.LogLevelDebug},
		{"info", "info", log.LogLevelInfo},
		{"INFO", "INFO", log.LogLevelInfo},
		{"warn", "warn", log.LogLevelWarning},
		{"WARNING", "WARNING", log.LogLevelWarning},
		{"warning", "warning", log.LogLevelWarning},
		{"error", "error", log.LogLevelError},
		{"ERROR", "ERROR", log.LogLevelError},
		{"critical", "critical", log.LogLevelCritical},
		{"CRITICAL", "CRITICAL", log.LogLevelCritical},
		{"invalid", "invalid", log.LogLevelInfo}, // 默认为 info
		{"empty", "", log.LogLevelInfo}, // 默认为 info
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := log.ParseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLevel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestLogModule 测试日志模块
func TestLogModule(t *testing.T) {
	tests := []struct {
		name   string
		module log.LogModule
		string string
	}{
		{"API模块", log.LogModuleAPI, "api"},
		{"Service模块", log.LogModuleService, "service"},
		{"Repository模块", log.LogModuleRepository, "repository"},
		{"Middleware模块", log.LogModuleMiddleware, "middleware"},
		{"Config模块", log.LogModuleConfig, "config"},
		{"Database模块", log.LogModuleDatabase, "database"},
		{"Cache模块", log.LogModuleCache, "cache"},
		{"External模块", log.LogModuleExternal, "external"},
		{"System模块", log.LogModuleSystem, "system"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.module.String() != tt.string {
				t.Errorf("LogModule.String() = %v, want %v", tt.module.String(), tt.string)
			}
		})
	}
}
