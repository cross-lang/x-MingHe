package validators

import (
	"reflect"
	"testing"
	"portal/internal/pkg/ginx"
	"github.com/go-playground/validator/v10"
)

// TestValidatePhone 测试手机号验证
func TestValidatePhone(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Phone string `validate:"phone"`
	}

	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"有效手机号", "13812345678", false},
		{"有效手机号2", "15987654321", false},
		{"有效手机号3", "18600001111", false},
		{"固定电话", "010-12345678", false},
		{"固定电话无连字符", "01012345678", false},
		{"空字符串", "", true},
		{"过短手机号", "138123456", true},
		{"过长手机号", "138123456789", true},
		{"错误前缀", "12812345678", true},
		{"非数字字符", "abc12345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Phone: tt.phone}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateMobile 测试手机号验证
func TestValidateMobile(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Mobile string `validate:"mobile"`
	}

	tests := []struct {
		name    string
		mobile  string
		wantErr bool
	}{
		{"有效手机号", "13812345678", false},
		{"有效手机号2", "15987654321", false},
		{"有效手机号3", "18600001111", false},
		{"固定电话", "010-12345678", true},
		{"空字符串", "", true},
		{"过短手机号", "138123456", true},
		{"过长手机号", "138123456789", true},
		{"错误前缀", "12812345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Mobile: tt.mobile}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMobile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateIDCard 测试身份证号验证
func TestValidateIDCard(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		IDCard string `validate:"idcard"`
	}

	tests := []struct {
		name    string
		idCard  string
		wantErr bool
	}{
		{"有效18位身份证", "110101199001011234", false},
		{"有效18位身份证带X", "11010119900101123X", false},
		{"有效18位身份证带小写x", "11010119900101123x", false},
		{"有效15位身份证", "110101900101123", false},
		{"空字符串", "", true},
		{"过短", "12345", true},
		{"过长", "110101199001011234567890", true},
		{"非数字", "abcdefghijklmnopqr", true},
		{"带空格", "110101 199001011234", false}, // 会移除空格
		{"带连字符", "110101-19900101-1234", false}, // 会移除连字符
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{IDCard: tt.idCard}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIDCard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidatePassword 测试密码强度验证
func TestValidatePassword(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Password string `validate:"password"`
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"有效密码", "Password123", false},
		{"有效密码2", "abc12345", false},
		{"有效密码3", "ABCD1234", false},
		{"有效密码4", "aB1cD2eF3", false},
		{"过短密码", "Ab1", true},
		{"纯字母", "abcdefgh", true},
		{"纯数字", "12345678", true},
		{"空字符串", "", true},
		{"只有7位", "Abc1234", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Password: tt.password}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateDate 测试日期验证
func TestValidateDate(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Date string `validate:"date"`
	}

	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"有效日期", "2023-01-01", false},
		{"有效日期2", "2023-12-31", false},
		{"有效日期3", "2000-02-29", false}, // 闰年
		{"空字符串", "", true},
		{"错误格式", "2023/01/01", true},
		{"错误格式2", "01-01-2023", true},
		{"无效日期", "2023-13-01", true},
		{"无效日期2", "2023-02-30", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Date: tt.date}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateDateTime 测试日期时间验证
func TestValidateDateTime(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		DateTime string `validate:"datetime"`
	}

	tests := []struct {
		name     string
		datetime string
		wantErr  bool
	}{
		{"有效日期时间", "2023-01-01 12:00:00", false},
		{"有效日期时间2", "2023-12-31 23:59:59", false},
		{"空字符串", "", true},
		{"错误格式", "2023/01/01 12:00:00", true},
		{"错误格式2", "2023-01-01T12:00:00", true},
		{"无效时间", "2023-01-01 25:00:00", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{DateTime: tt.datetime}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDateTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateEmailStrict 测试严格邮箱验证
func TestValidateEmailStrict(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Email string `validate:"email_strict"`
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"有效邮箱", "test@example.com", false},
		{"有效邮箱2", "user.name@domain.co.uk", false},
		{"有效邮箱3", "user+tag@example.com", false},
		{"空字符串", "", true},
		{"缺少@", "testexample.com", true},
		{"缺少域名", "test@", true},
		{"缺少用户名", "@example.com", true},
		{"无效字符", "test@exa mple.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Email: tt.email}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmailStrict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateUsername 测试用户名验证
func TestValidateUsername(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Username string `validate:"username"`
	}

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"有效用户名", "user123", false},
		{"有效用户名2", "test_user", false},
		{"有效用户名3", "admin", false},
		{"空字符串", "", true},
		{"过短", "abc", true},
		{"过长", "this_is_a_very_long_username_that_exceeds_twenty_characters", true},
		{"包含特殊字符", "user@123", true},
		{"包含空格", "user 123", true},
		{"包含中文", "用户123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Username: tt.username}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateURL 测试URL验证
func TestValidateURL(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		URL string `validate:"url"`
	}

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"有效HTTP URL", "http://example.com", false},
		{"有效HTTPS URL", "https://example.com", false},
		{"有效URL带路径", "https://example.com/path/to/resource", false},
		{"有效URL带参数", "https://example.com?param=value", false},
		{"空字符串", "", true},
		{"缺少协议", "example.com", true},
		{"无效协议", "ftp://example.com", true},
		{"包含空格", "http://example .com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{URL: tt.url}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateStringRange 测试字符串长度范围验证
func TestValidateStringRange(t *testing.T) {
	validate := validator.New()
	if err := ginx.InitCustomValidators(validate); err != nil {
		t.Fatalf("初始化验证器失败: %v", err)
	}

	type TestStruct struct {
		Name string `validate:"string_range=2,10"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"有效长度", "张三", false},
		{"有效长度2", "hello", false},
		{"最小边界", "李", false},
		{"最大边界", "1234567890", false},
		{"过短", "A", true},
		{"过长", "这个名字太长了超过了限制", true},
		{"空字符串", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Name: tt.value}
			err := validate.Struct(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateStringRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidationMessages 测试验证消息
func TestValidationMessages(t *testing.T) {
	messages := ginx.ValidationMessages

	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{"手机号消息", ginx.TagPhone, "请输入有效的电话号码"},
		{"密码消息", ginx.TagPassword, "密码至少8位，且必须包含字母和数字"},
		{"身份证消息", ginx.TagIDCard, "请输入有效的身份证号码"},
		{"日期消息", ginx.TagDate, "请输入有效的日期格式（YYYY-MM-DD）"},
		{"用户名消息", ginx.TagUsername, "用户名只能包含字母、数字、下划线，长度4-20位"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if msg, ok := messages[tt.tag]; !ok {
				t.Errorf("验证器标签 %s 的消息不存在", tt.tag)
			} else if msg != tt.expected {
				t.Errorf("验证器标签 %s 的消息不匹配，期望 %s，实际 %s", tt.tag, tt.expected, msg)
			}
		})
	}
}

// TestValidatorTags 测试验证器标签常量
func TestValidatorTags(t *testing.T) {
	tags := []string{
		ginx.TagPhone,
		ginx.TagMobile,
		ginx.TagIDCard,
		ginx.TagPassword,
		ginx.TagDate,
		ginx.TagDateTime,
		ginx.TagFuture,
		ginx.TagPast,
		ginx.TagEnum,
		ginx.TagStringRange,
		ginx.TagEmailStrict,
		ginx.TagUsername,
		ginx.TagURL,
		ginx.TagGteField,
		ginx.TagLteField,
	}

	for _, tag := range tags {
		if tag == "" {
			t.Errorf("验证器标签为空")
		}
	}
}
