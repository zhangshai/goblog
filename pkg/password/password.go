package password

import (
	"goblog/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// Hash 使用 bcrypt 对密码进行加密
func Hash(str string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	logger.LogError(err)
	return string(bytes)
}

// CheckHash 对比明文密码和数据库的哈希值
func CheckHash(hash, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	logger.LogError(err)
	return err == nil
}

func IsHashed(str string) bool {

	// bcrypt 加密后的长度等于 60
	return len(str) == 60
}
