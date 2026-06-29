package timex

import (
	"strings"
	"time"
)

func ParseChineseDateTime(s string) (time.Time, error) {
	// 替换中文上午/下午为英文
	s = strings.ReplaceAll(s, "上午", "AM")
	s = strings.ReplaceAll(s, "下午", "PM")

	// 定义 layout（注意：12小时制用 03，24小时制用 15）
	layout := "2006年1月2日 PM03:04"

	// 建议指定时区（如北京时间）
	beijing, _ := time.LoadLocation("Asia/Shanghai")
	return time.ParseInLocation(layout, s, beijing)
}

// TimestampToTime 将秒级时间戳转换为time.Time类型
func TimestampToTime(timestamp int64) time.Time {
	// 第二个参数为纳秒偏移量，这里传入0表示不添加额外偏移
	return time.Unix(timestamp, 0)
}

// TimestampMsToTime 将毫秒级时间戳转换为time.Time类型
func TimestampMsToTime(timestamp int64) time.Time {
	return time.UnixMilli(timestamp)
}
