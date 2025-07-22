package jwtutil

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

// GenerateToken 生成 JWT Token
// 参数说明：
//   - userId: 用户ID
//   - userType: 用户类型（0: 管理员, 1: 教师, 2: 学生）
//   - secret: JWT签名密钥
//   - expireDuration: Token的有效期，例如：24 * time.Hour
func GenerateToken(userId int, userType int, secret string, expireDuration time.Duration) (string, error) {
	if secret == "" {
		return "", errors.New("JWT密钥未配置")
	}

	// 定义 JWT Claims
	claims := jwt.MapClaims{
		"userId":   userId,
		"userType": userType, // 新增用户类型字段
		"exp":      time.Now().Add(expireDuration).Unix(),
	}

	// 使用HS256签名方法生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ParseToken 解析JWT Token并返回userId 和userType
func ParseToken(tokenString string, secret string) (int64, int, error) {
	if secret == "" {
		return 0, 0, errors.New("JWT密钥未配置")
	}

	// 去掉"Bearer "前缀（如果存在）
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 解析Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保 Token 使用的是 HS256 签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token签名方法无效")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, 0, err
	}

	// 从Claims中提取userId和userType
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdFloat, userIdOk := claims["userId"].(float64)
		userTypeFloat, userTypeOk := claims["userType"].(float64)

		if userIdOk && userTypeOk {
			return int64(userIdFloat), int(userTypeFloat), nil
		}
	}

	return 0, 0, errors.New("无效的token")
}
