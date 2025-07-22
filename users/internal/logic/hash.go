package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// GenerateSalt  生成四位随机盐
func GenerateSalt() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(9999)
}

// HashPassword 使用盐进行MD5加密
func HashPassword(password string, salt int64) string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d%s", salt, password)))
	return hex.EncodeToString(hash.Sum(nil))
}
