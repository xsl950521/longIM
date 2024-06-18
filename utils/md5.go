package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// 转小写
func Md5Encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 转大写
func MD5Encode(str string) string {
	return strings.ToUpper(Md5Encode(str))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, password string) bool {
	md := Md5Encode(plainpwd + salt)
	fmt.Println("plainpwd=", plainpwd, "   salt=", salt, "   md=", md, "         password=", password)
	return Md5Encode(plainpwd+salt) == password
}
