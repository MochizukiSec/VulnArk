package utils

import (
	"strconv"
)

// ParseInt 将字符串转换为整数，如果转换失败则返回默认值
func ParseInt(str string, defaultValue int) int {
	if value, err := strconv.Atoi(str); err == nil {
		return value
	}
	return defaultValue
}
