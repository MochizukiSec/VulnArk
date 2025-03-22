package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// 在出错的情况下返回固定值
		return "ERROR_GENERATING_RANDOM_STRING"
	}
	return hex.EncodeToString(bytes)
}
