package ginx

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// InitCustomValidators 初始化自定义验证器
func InitCustomValidators(v *validator.Validate) error {
	// 注册自定义验证器
	if err := v.RegisterValidation("phone", validatePhone); err != nil {
		return fmt.Errorf("注册 phone 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("mobile", validateMobile); err != nil {
		return fmt.Errorf("注册 mobile 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("idcard", validateIDCard); err != nil {
		return fmt.Errorf("注册 idcard 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("password", validatePassword); err != nil {
		return fmt.Errorf("注册 password 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("date", validateDate); err != nil {
		return fmt.Errorf("注册 date 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("datetime", validateDateTime); err != nil {
		return fmt.Errorf("注册 datetime 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("future", validateFuture); err != nil {
		return fmt.Errorf("注册 future 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("past", validatePast); err != nil {
		return fmt.Errorf("注册 past 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("enum", validateEnum); err != nil {
		return fmt.Errorf("注册 enum 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("string_range", validateStringRange); err != nil {
		return fmt.Errorf("注册 string_range 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("email_strict", validateEmailStrict); err != nil {
		return fmt.Errorf("注册 email_strict 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("username", validateUsername); err != nil {
		return fmt.Errorf("注册 username 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("url", validateURL); err != nil {
		return fmt.Errorf("注册 url 验证器失败: %w", err)
	}

	// 注册自定义验证器函数
	if err := v.RegisterValidation("gte_field", validateGteField); err != nil {
		return fmt.Errorf("注册 gte_field 验证器失败: %w", err)
	}
	if err := v.RegisterValidation("lte_field", validateLteField); err != nil {
		return fmt.Errorf("注册 lte_field 验证器失败: %w", err)
	}

	return nil
}

// validatePhone 验证电话号码（支持固定电话和手机号）
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return false
	}

	// 匹配手机号或固定电话
	mobileRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	landlineRegex := regexp.MustCompile(`^0\d{2,3}-?\d{7,8}$`)

	return mobileRegex.MatchString(phone) || landlineRegex.MatchString(phone)
}

// validateMobile 验证手机号（中国大陆）
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if mobile == "" {
		return false
	}

	// 中国大陆手机号：1开头，第二位3-9，共11位
	mobileRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return mobileRegex.MatchString(mobile)
}

// validateIDCard 验证身份证号（中国大陆）
func validateIDCard(fl validator.FieldLevel) bool {
	idCard := fl.Field().String()
	if idCard == "" {
		return false
	}

	// 移除所有空格
	idCard = strings.ReplaceAll(idCard, " ", "")
	idCard = strings.ReplaceAll(idCard, "-", "")

	// 15位身份证
	if len(idCard) == 15 {
		regex := regexp.MustCompile(`^\d{15}$`)
		return regex.MatchString(idCard)
	}

	// 18位身份证
	if len(idCard) == 18 {
		regex := regexp.MustCompile(`^\d{17}[\dXx]$`)
		return regex.MatchString(idCard)
	}

	return false
}

// validatePassword 验证密码强度
// 至少8位，包含字母和数字
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasLetter && hasDigit
}

// validateDate 验证日期格式（YYYY-MM-DD）
func validateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return false
	}

	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// validateDateTime 验证日期时间格式（YYYY-MM-DD HH:mm:ss）
func validateDateTime(fl validator.FieldLevel) bool {
	dateTimeStr := fl.Field().String()
	if dateTimeStr == "" {
		return false
	}

	_, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	return err == nil
}

// validateFuture 验证日期是否在未来
func validateFuture(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return false
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	return date.After(time.Now())
}

// validatePast 验证日期是否在过去
func validatePast(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return false
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	return date.Before(time.Now())
}

// validateEnum 验证枚举值
// 使用方法：`enum=value1,value2,value3`
func validateEnum(fl validator.FieldLevel) bool {
	param := fl.Param()
	field := fl.Field()

	// 获取允许的枚举值
	allowedValues := strings.Split(param, ",")

	// 根据字段类型进行验证
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		for _, allowed := range allowedValues {
			if value == allowed {
				return true
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := field.Int()
		for _, allowed := range allowedValues {
			if allowedInt, err := strconv.ParseInt(allowed, 10, 64); err == nil && value == allowedInt {
				return true
			}
		}
	}

	return false
}

// validateStringRange 验证字符串长度范围
// 使用方法：`string_range=min,max`
func validateStringRange(fl validator.FieldLevel) bool {
	param := fl.Param()
	field := fl.Field().String()

	parts := strings.Split(param, ",")
	if len(parts) != 2 {
		return false
	}

	minLen, err1 := strconv.Atoi(parts[0])
	maxLen, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return false
	}

	length := len([]rune(field))
	return length >= minLen && length <= maxLen
}

// validateEmailStrict 严格验证邮箱格式
func validateEmailStrict(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return false
	}

	// 更严格的邮箱正则表达式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validateUsername 验证用户名
// 只能包含字母、数字、下划线，4-20位
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return false
	}

	// 长度限制
	if len(username) < 4 || len(username) > 20 {
		return false
	}

	// 只能包含字母、数字、下划线
	regex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return regex.MatchString(username)
}

// validateURL 验证 URL 格式
func validateURL(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	if urlStr == "" {
		return false
	}

	// 简单的 URL 验证
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(urlStr)
}

// validateGteField 验证字段是否大于等于另一个字段
// 使用方法：`gte_field=other_field`
func validateGteField(fl validator.FieldLevel) bool {
	otherFieldName := fl.Param()
	otherField := fl.Top().FieldByName(otherFieldName)

	if !otherField.IsValid() {
		return false
	}

	// 比较两个字段
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() >= otherField.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Field().Uint() >= otherField.Uint()
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() >= otherField.Float()
	}

	return false
}

// validateLteField 验证字段是否小于等于另一个字段
// 使用方法：`lte_field=other_field`
func validateLteField(fl validator.FieldLevel) bool {
	otherFieldName := fl.Param()
	otherField := fl.Top().FieldByName(otherFieldName)

	if !otherField.IsValid() {
		return false
	}

	// 比较两个字段
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fl.Field().Int() <= otherField.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fl.Field().Uint() <= otherField.Uint()
	case reflect.Float32, reflect.Float64:
		return fl.Field().Float() <= otherField.Float()
	}

	return false
}

// ValidatorTags 验证器标签说明
const (
	// TagPhone 电话号码验证（支持固定电话和手机号）
	TagPhone = "phone"
	// TagMobile 手机号验证（中国大陆）
	TagMobile = "mobile"
	// TagIDCard 身份证号验证（中国大陆）
	TagIDCard = "idcard"
	// TagPassword 密码强度验证（至少8位，包含字母和数字）
	TagPassword = "password"
	// TagDate 日期格式验证（YYYY-MM-DD）
	TagDate = "date"
	// TagDateTime 日期时间格式验证（YYYY-MM-DD HH:mm:ss）
	TagDateTime = "datetime"
	// TagFuture 日期是否在未来
	TagFuture = "future"
	// TagPast 日期是否在过去
	TagPast = "past"
	// TagEnum 枚举值验证
	TagEnum = "enum"
	// TagStringRange 字符串长度范围
	TagStringRange = "string_range"
	// TagEmailStrict 严格邮箱验证
	TagEmailStrict = "email_strict"
	// TagUsername 用户名验证
	TagUsername = "username"
	// TagURL URL 验证
	TagURL = "url"
	// TagGteField 大于等于另一个字段
	TagGteField = "gte_field"
	// TagLteField 小于等于另一个字段
	TagLteField = "lte_field"
)

// ValidationMessages 自定义验证错误消息
var ValidationMessages = map[string]string{
	TagPhone:        "请输入有效的电话号码",
	TagMobile:       "请输入有效的手机号码",
	TagIDCard:       "请输入有效的身份证号码",
	TagPassword:     "密码至少8位，且必须包含字母和数字",
	TagDate:         "请输入有效的日期格式（YYYY-MM-DD）",
	TagDateTime:     "请输入有效的日期时间格式（YYYY-MM-DD HH:mm:ss）",
	TagFuture:       "日期必须在未来",
	TagPast:         "日期必须在过去",
	TagEnum:         "请选择有效的值",
	TagStringRange:  "字符串长度不符合要求",
	TagEmailStrict:  "请输入有效的邮箱地址",
	TagUsername:     "用户名只能包含字母、数字、下划线，长度4-20位",
	TagURL:          "请输入有效的URL地址",
	TagGteField:     "值必须大于等于指定字段",
	TagLteField:     "值必须小于等于指定字段",
}
