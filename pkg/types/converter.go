package types

import (
	"goblog/pkg/logger"
	"math/rand"
	"strconv"
	"time"
)

// Int64ToString 将 int64 转换为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func StringToUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logger.LogError(err)
	}
	return i
}

func StringToInt(str string) int {
	_str, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logger.LogError(err)
	}
	return int(_str)
}

// Int64ToString 将 int64 转换为 string
func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func RandStr(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
