package utils

import "fmt"

func GetActivityTypeName(t int32) string {
	switch t {
	case 1:
		return "课程"
	case 2:
		return "活动"
	default:
		return fmt.Sprintf("未知(%d)", t)
	}
}

func GetActivitySignUpStatusName(s int32) string {
	switch s {
	case 1:
		return "报名成功"
	case 2:
		return "报名取消"
	case 3:
		return "报名失败"
	default:
		return fmt.Sprintf("未知(%d)", s)
	}
}
