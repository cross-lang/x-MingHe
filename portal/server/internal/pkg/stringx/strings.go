package stringx

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Split(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return strings.Split(s, sep)
}

func StrLike(str string) string {
	str = strings.ReplaceAll(str, "%", "\\%")
	str = strings.ReplaceAll(str, "_", "\\_")

	return "%" + str + "%"
}

func GetTradeNo(uniq string) string {
	randData := rand.Intn(100000)
	return fmt.Sprintf("TR%s%s%d", uniq, time.Now().Format("20060102150405"), randData)
}
