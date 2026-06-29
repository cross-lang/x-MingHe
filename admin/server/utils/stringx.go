package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	rand2 "math/rand"
	"strconv"
	"strings"
	"time"
)

func Split(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return strings.Split(s, sep)
}

const (
	CharSetV1 = "0123456789"
	CharSetV2 = "abcdefghijklmnopqrstuvwxyz"
	CharSetV3 = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	CharSetV4 = CharSetV1 + CharSetV2
	CharSetV5 = CharSetV1 + CharSetV3
	CharSetV6 = CharSetV1 + CharSetV2 + CharSetV3
)

// GenerateRandString 根据字符集合和长度生成验证码
func GenerateRandString(charset string, length int) string {
	if length <= 0 {
		return ""
	}
	if charset == "" {
		charset = CharSetV6
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

// GenerateTradeNo 生成交易流水号
func GenerateTradeNo(id int) string {
	randData := rand2.Intn(100000)
	return fmt.Sprintf("TR%s%s%d", strconv.Itoa(id), time.Now().Format("20060102150405"), randData)
}
