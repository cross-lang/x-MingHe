package log

import (
	"regexp"
	"strings"
)

// MaskType 脱敏类型
type MaskType string

const (
	// MaskTypePhone 手机号脱敏
	MaskTypePhone MaskType = "phone"
	// MaskTypeEmail 邮箱脱敏
	MaskTypeEmail MaskType = "email"
	// MaskTypeIDCard 身份证脱敏
	MaskTypeIDCard MaskType = "idcard"
	// MaskTypeBankCard 银行卡脱敏
	MaskTypeBankCard MaskType = "bankcard"
	// MaskTypePassword 密码脱敏
	MaskTypePassword MaskType = "password"
	// MaskTypeName 姓名（部分）脱敏
	MaskTypeName MaskType = "name"
	// MaskTypeAddress 地址脱敏
	MaskTypeAddress MaskType = "address"
)

// Mask 脱敏函数
func Mask(data string, maskType MaskType) string {
	if data == "" {
		return data
	}

	switch maskType {
	case MaskTypePhone:
		return maskPhone(data)
	case MaskTypeEmail:
		return maskEmail(data)
	case MaskTypeIDCard:
		return maskIDCard(data)
	case MaskTypeBankCard:
		return maskBankCard(data)
	case MaskTypePassword:
		return "******"
	case MaskTypeName:
		return maskName(data)
	case MaskTypeAddress:
		return maskAddress(data)
	default:
		return data
	}
}

// maskPhone 手机号脱敏：保留前3位和后4位，中间用****代替
// 示例：13812345678 -> 138****5678
func maskPhone(phone string) string {
	if len(phone) < 7 {
		return "****"
	}

	// 移除所有非数字字符
	re := regexp.MustCompile(`[^\d]`)
	phone = re.ReplaceAllString(phone, "")

	if len(phone) != 11 {
		return phone[:3] + "****" + phone[len(phone)-4:]
	}

	return phone[:3] + "****" + phone[7:]
}

// maskEmail 邮箱脱敏：保留邮箱前缀第一个字符和@及域名
// 示例：example@test.com -> e****@test.com
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "****@****.com"
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 1 {
		return "****@" + domain
	}

	return username[:1] + "****@" + domain
}

// maskIDCard 身份证号脱敏：保留前6位和后4位
// 示例：110101199001011234 -> 110101********1234
func maskIDCard(idCard string) string {
	if len(idCard) < 10 {
		return "****"
	}

	// 移除所有非字母数字字符
	re := regexp.MustCompile(`[^\da-zA-Z]`)
	idCard = re.ReplaceAllString(idCard, "")

	length := len(idCard)
	if length == 15 {
		return idCard[:6] + "*******" + idCard[12:]
	}
	if length == 18 {
		return idCard[:6] + "********" + idCard[14:]
	}

	return idCard[:6] + "****" + idCard[length-4:]
}

// maskBankCard 银行卡号脱敏：保留前6位和后4位
// 示例：6225880212345678 -> 622588****5678
func maskBankCard(cardNum string) string {
	// 移除所有非数字字符
	re := regexp.MustCompile(`[^\d]`)
	cardNum = re.ReplaceAllString(cardNum, "")

	if len(cardNum) < 10 {
		return "****"
	}

	return cardNum[:6] + "****" + cardNum[len(cardNum)-4:]
}

// maskName 姓名脱敏：保留姓氏，名字用*代替
// 示例：张三 -> 张*，欧阳娜娜 -> 欧阳**
func maskName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return "*"
	}

	// 保留第一个字符（姓氏）
	result := string(runes[:1])

	// 姓氏可能是两个字符（如欧阳）
	if len(runes) > 2 && isChineseSurname(string(runes[:2])) {
		result = string(runes[:2])
	}

	// 剩余字符用*代替
	for i := len([]rune(result)); i < len(runes); i++ {
		result += "*"
	}

	return result
}

// maskAddress 地址脱敏：保留前3个字符和后3个字符
// 示例：北京市朝阳区某某路123号 -> 北京市...23号
func maskAddress(address string) string {
	if len(address) < 10 {
		return "****"
	}

	runes := []rune(address)
	return string(runes[:3]) + "..." + string(runes[len(runes)-3:])
}

// isChineseSurname 判断是否是常见的中文复姓
func isChineseSurname(surname string) bool {
	commonCompoundSurnames := []string{
		"欧阳", "太史", "端木", "上官", "司马", "东方", "独孤", "南宫", "万俟",
		"闻人", "夏侯", "诸葛", "尉迟", "公羊", "赫连", "澹台", "皇甫", "宗政",
		"濮阳", "公西", "太叔", "申屠", "公孙", "慕容", "仲孙", "钟离", "长孙",
		"宇文", "司徒", "鲜于", "司空", "闾丘", "子车", "亓官", "司寇", "巫马",
		"公西", "颛孙", "壤驷", "公良", "拓跋", "夹谷", "宰父", "谷梁", "晋楚",
		"法汝", "邾", "淳于", "单于", "令狐", "郭", "夏", "侯", "刁", "袁",
	}

	for _, s := range commonCompoundSurnames {
		if surname == s {
			return true
		}
	}
	return false
}

// MaskField 脱敏字段
// 根据字段名自动判断脱敏类型并返回脱敏后的值
func MaskField(fieldName, value string) string {
	if value == "" {
		return value
	}

	fieldName = strings.ToLower(fieldName)

	// 密码相关字段
	if strings.Contains(fieldName, "password") ||
		strings.Contains(fieldName, "passwd") ||
		strings.Contains(fieldName, "pwd") {
		return Mask(value, MaskTypePassword)
	}

	// 手机号相关字段
	if strings.Contains(fieldName, "phone") ||
		strings.Contains(fieldName, "mobile") ||
		strings.Contains(fieldName, "tel") {
		return Mask(value, MaskTypePhone)
	}

	// 邮箱相关字段
	if strings.Contains(fieldName, "email") ||
		strings.Contains(fieldName, "mail") {
		return Mask(value, MaskTypeEmail)
	}

	// 身份证相关字段
	if strings.Contains(fieldName, "idcard") ||
		strings.Contains(fieldName, "id_card") ||
		strings.Contains(fieldName, "identity") {
		return Mask(value, MaskTypeIDCard)
	}

	// 银行卡相关字段
	if strings.Contains(fieldName, "bankcard") ||
		strings.Contains(fieldName, "card") ||
		strings.Contains(fieldName, "account") {
		return Mask(value, MaskTypeBankCard)
	}

	// 姓名相关字段
	if strings.Contains(fieldName, "name") {
		return Mask(value, MaskTypeName)
	}

	// 地址相关字段
	if strings.Contains(fieldName, "address") ||
		strings.Contains(fieldName, "addr") {
		return Mask(value, MaskTypeAddress)
	}

	return value
}

// MaskMap 脱敏 map 中的敏感字段
func MaskMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			result[key] = MaskField(key, v)
		case map[string]interface{}:
			result[key] = MaskMap(v)
		default:
			result[key] = value
		}
	}
	return result
}
