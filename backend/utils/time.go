package utils

import (
	"time"
)

// 中国时区
var CSTZone = time.FixedZone("CST", 8*3600) // 东八区

// FormatTimeCST 将时间格式化为中国标准时间格式
func FormatTimeCST(t time.Time) string {
	// 将时间转换为北京时区
	cstTime := t.In(CSTZone)
	// 格式化为易读的格式
	return cstTime.Format("2006-01-02 15:04:05")
}

// NowCST 获取当前北京时间
func NowCST() time.Time {
	return time.Now().In(CSTZone)
}
