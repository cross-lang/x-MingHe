package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// MD5WithSalt 对密码进行MD5加盐哈希
// 返回格式：盐值$哈希值
func MD5WithSalt(password string, salt string) string {
	// 组合密码和盐值，使用$分隔符
	combined := fmt.Sprintf("%s$%s", password, salt)
	// 计算MD5哈希
	hash := md5.Sum([]byte(combined))
	// 将哈希值转换为十六进制字符串
	hashHex := hex.EncodeToString(hash[:])
	// 返回盐值$哈希值的格式
	return hashHex
}
