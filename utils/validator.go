package utils

import (
	"regexp"
)

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// IsPhone 验证手机号格式（中国大陆）
func IsPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// IsUsername 验证用户名格式（3-20位字母、数字、下划线）
func IsUsername(username string) bool {
	pattern := `^[a-zA-Z0-9_]{3,20}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}
